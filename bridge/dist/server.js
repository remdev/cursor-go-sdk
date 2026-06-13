import { createHash } from "node:crypto";
import http from "node:http";
import { homedir } from "node:os";
import { join, resolve } from "node:path";
import { SdkAgentService } from "@anysphere/proto/sdk/v1/sdk_agent_service_connect.js";
import { SdkBridgeControlService } from "@anysphere/proto/sdk/v1/sdk_bridge_control_service_connect.js";
import { SdkCursorService } from "@anysphere/proto/sdk/v1/sdk_cursor_service_connect.js";
import { connectNodeAdapter } from "@connectrpc/connect-node";
import { createBridgeAuthInterceptor, createBridgeAuthToken } from "./auth.js";
import { CURSOR_SDK_BRIDGE_VERSION } from "./constants.js";
import { CursorSdkBridgeRegistry } from "./registry.js";
import { createSdkErrorInterceptor } from "./sdk-error-interceptor.js";
import { createSdkAgentService, createSdkBridgeControlService, createSdkCursorService, } from "./sdk-service.js";
export async function startCursorSdkBridgeServer(options = {}) {
    var _a;
    var _b, _c, _d, _e, _f, _g;
    const host = (_b = options.host) !== null && _b !== void 0 ? _b : "127.0.0.1";
    if (!options.allowNonLoopbackHost && !isLoopbackHost(host)) {
        throw new Error(`Refusing to bind Cursor SDK bridge to non-loopback host "${host}" without allowNonLoopbackHost`);
    }
    const port = (_c = options.port) !== null && _c !== void 0 ? _c : 0;
    const authToken = (_d = options.authToken) !== null && _d !== void 0 ? _d : createBridgeAuthToken();
    const registry = (_e = options.registry) !== null && _e !== void 0 ? _e : new CursorSdkBridgeRegistry();
    // Resolved launch workspace. Used both for the advertised address below and as
    // the default `cwd` for local-store ops (create / list / get) so reads find
    // agents created here without a per-call `cwd`. See sdk-service.ts.
    const workspaceRef = resolve((_f = options.workspaceRef) !== null && _f !== void 0 ? _f : process.cwd());
    let closeBridge;
    const handler = connectNodeAdapter({
        routes: router => {
            mountSdkBridgeRoutes(router, {
                registry,
                workspaceRef,
                shutdown: async graceSeconds => {
                    const close = closeBridge;
                    const delayMs = Math.max(0, Math.trunc(graceSeconds)) * 1000;
                    const timeout = setTimeout(() => {
                        void (close === null || close === void 0 ? void 0 : close());
                    }, delayMs);
                    timeout.unref();
                },
            });
        },
        // Order matters: auth runs first so unauthenticated requests never reach
        // the error translator. Subsequent interceptors run in array order, so
        // errors thrown by the SDK funnel through the translator on the way out.
        interceptors: [
            createBridgeAuthInterceptor(authToken),
            createSdkErrorInterceptor(),
        ],
        requireConnectProtocolHeader: false,
    });
    const server = http.createServer(handler);
    await new Promise((resolve, reject) => {
        server.once("error", reject);
        server.listen(port, host, () => {
            server.off("error", reject);
            resolve();
        });
    });
    const serverAddress = server.address();
    if (serverAddress === null ||
        typeof serverAddress === "string" ||
        typeof serverAddress.port !== "number") {
        await closeServer(server);
        throw new Error("Cursor SDK bridge failed to determine listening address");
    }
    const stateRoot = (_g = options.stateRoot) !== null && _g !== void 0 ? _g : defaultStateRootForWorkspace(workspaceRef);
    const address = {
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
    (_a = options.writeReady) === null || _a === void 0 ? void 0 : _a.call(options, address);
    const handle = {
        address,
        close: async () => {
            let closeError;
            try {
                await closeServer(server);
            }
            catch (err) {
                closeError = err;
            }
            try {
                await registry.dispose();
            }
            catch (disposeError) {
                if (closeError) {
                    throw new AggregateError([closeError, disposeError], "Failed to close Cursor SDK bridge cleanly");
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
function mountSdkBridgeRoutes(router, options) {
    router.service(SdkAgentService, createSdkAgentService({
        registry: options.registry,
        workspaceRef: options.workspaceRef,
    }));
    router.service(SdkCursorService, createSdkCursorService());
    router.service(SdkBridgeControlService, createSdkBridgeControlService({ shutdown: options.shutdown }));
}
function closeServer(server) {
    return new Promise((resolve, reject) => {
        server.close(error => {
            if (error) {
                reject(error);
                return;
            }
            resolve();
        });
    });
}
function isLoopbackHost(host) {
    return host === "127.0.0.1" || host === "localhost" || host === "::1";
}
function formatBridgeServerUrl(host, port) {
    const urlHost = host.includes(":") && !host.startsWith("[") ? `[${host}]` : host;
    return new URL(`http://${urlHost}:${port}`).origin;
}
function defaultStateRootForWorkspace(workspaceRef) {
    const workspaceHash = createHash("sha256")
        .update(workspaceRef)
        .digest("hex")
        .slice(0, 24);
    return join(homedir(), ".cursor", "sdk-agent-store", workspaceHash);
}
export const startBridgeServer = startCursorSdkBridgeServer;
