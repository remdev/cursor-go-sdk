import { createHash } from "node:crypto";
import http from "node:http";
import { homedir } from "node:os";
import { join, resolve } from "node:path";
import { SdkAgentService } from "../gen/ts/sdk/v1/sdk_agent_service_connect.js";
import { SdkBridgeControlService } from "../gen/ts/sdk/v1/sdk_bridge_control_service_connect.js";
import { SdkCursorService } from "../gen/ts/sdk/v1/sdk_cursor_service_connect.js";
import { connectNodeAdapter } from "@connectrpc/connect-node";
import type { ConnectRouter } from "@connectrpc/connect";
import { createBridgeAuthInterceptor, createBridgeAuthToken } from "./auth.js";
import { CURSOR_SDK_BRIDGE_VERSION } from "./constants.js";
import { CursorSdkBridgeRegistry } from "./registry.js";
import { createSdkErrorInterceptor } from "./sdk-error-interceptor.js";
import {
    createSdkAgentService,
    createSdkBridgeControlService,
    createSdkCursorService,
} from "./sdk-service.js";
import type {
    BridgeServerAddress,
    MountSdkBridgeRoutesOptions,
    StartCursorSdkBridgeServerOptions,
} from "./types.js";

export interface BridgeServerHandle {
    address: BridgeServerAddress;
    close: () => Promise<void>;
}

export async function startCursorSdkBridgeServer(
    options: StartCursorSdkBridgeServerOptions = {},
): Promise<BridgeServerHandle> {
    const host = options.host ?? "127.0.0.1";
    if (!options.allowNonLoopbackHost && !isLoopbackHost(host)) {
        throw new Error(
            `Refusing to bind Cursor SDK bridge to non-loopback host "${host}" without allowNonLoopbackHost`,
        );
    }
    const port = options.port ?? 0;
    const authToken = options.authToken ?? createBridgeAuthToken();
    const registry = options.registry ?? new CursorSdkBridgeRegistry();
    // Resolved launch workspace. Used both for the advertised address below and as
    // the default `cwd` for local-store ops (create / list / get) so reads find
    // agents created here without a per-call `cwd`. See sdk-service.ts.
    const workspaceRef = resolve(options.workspaceRef ?? process.cwd());
    let closeBridge: (() => Promise<void>) | undefined;
    const handler = connectNodeAdapter({
        routes: (router) => {
            mountSdkBridgeRoutes(router, {
                registry,
                workspaceRef,
                shutdown: async (graceSeconds) => {
                    const close = closeBridge;
                    const delayMs = Math.max(0, Math.trunc(graceSeconds)) * 1000;
                    const timeout = setTimeout(() => {
                        void close?.();
                    }, delayMs);
                    timeout.unref();
                },
            });
        },
        interceptors: [
            createBridgeAuthInterceptor(authToken),
            createSdkErrorInterceptor(),
        ],
        requireConnectProtocolHeader: false,
    });
    const server = http.createServer(handler);
    await new Promise<void>((resolveListen, reject) => {
        server.once("error", reject);
        server.listen(port, host, () => {
            server.off("error", reject);
            resolveListen();
        });
    });
    const serverAddress = server.address();
    if (
        serverAddress === null ||
        typeof serverAddress === "string" ||
        typeof serverAddress.port !== "number"
    ) {
        await closeServer(server);
        throw new Error("Cursor SDK bridge failed to determine listening address");
    }
    const stateRoot = options.stateRoot ?? defaultStateRootForWorkspace(workspaceRef);
    const address: BridgeServerAddress = {
        schemaVersion: 1,
        serverVersion: CURSOR_SDK_BRIDGE_VERSION,
        pid: process.pid,
        transport: "tcp",
        protocol: "connect",
        host,
        port: serverAddress.port,
        url: formatBridgeServerUrl(host, serverAddress.port),
        authToken,
        workspaceRef,
        stateRoot,
        maxConcurrentAgents: options.maxConcurrentAgents,
        maxMessageBytes: options.maxMessageBytes,
    };
    options.writeReady?.(address);
    const handle: BridgeServerHandle = {
        address,
        close: async () => {
            let closeError: unknown;
            try {
                await closeServer(server);
            } catch (err) {
                closeError = err;
            }
            try {
                await registry.dispose();
            } catch (disposeError) {
                if (closeError) {
                    throw new AggregateError(
                        [closeError, disposeError],
                        "Failed to close Cursor SDK bridge cleanly",
                    );
                }
                throw disposeError;
            }
            if (closeError) {
                throw closeError;
            }
        },
    };
    closeBridge = handle.close;
    return handle;
}

function mountSdkBridgeRoutes(router: ConnectRouter, options: MountSdkBridgeRoutesOptions): void {
    router.service(
        SdkAgentService,
        createSdkAgentService({
            registry: options.registry,
            workspaceRef: options.workspaceRef,
        }),
    );
    router.service(SdkCursorService, createSdkCursorService());
    router.service(
        SdkBridgeControlService,
        createSdkBridgeControlService({ shutdown: options.shutdown }),
    );
}

function closeServer(server: http.Server): Promise<void> {
    return new Promise((resolveClose, reject) => {
        server.close((error) => {
            if (error) {
                reject(error);
                return;
            }
            resolveClose();
        });
    });
}

function isLoopbackHost(host: string): boolean {
    return host === "127.0.0.1" || host === "localhost" || host === "::1";
}

function formatBridgeServerUrl(host: string, port: number): string {
    const urlHost = host.includes(":") && !host.startsWith("[") ? `[${host}]` : host;
    return new URL(`http://${urlHost}:${port}`).origin;
}

function defaultStateRootForWorkspace(workspaceRef: string): string {
    const workspaceHash = createHash("sha256")
        .update(workspaceRef)
        .digest("hex")
        .slice(0, 24);
    return join(homedir(), ".cursor", "sdk-agent-store", workspaceHash);
}

export const startBridgeServer = startCursorSdkBridgeServer;
