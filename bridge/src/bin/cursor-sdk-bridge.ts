#!/usr/bin/env node
import { mkdtempSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { LocalAgentStoreConfig } from "../../gen/ts/sdk/v1/sdk_messages_pb.js";
import { configureCursorSdk } from "@cursor/sdk";
import { createLocalAgentStoreFromProtoConfig } from "../bridge-local-agent-store.js";
import { startCursorSdkBridgeServer } from "../server.js";
import { readBridgeToolCallbackConfigFromEnv, setBridgeToolCallbackConfig } from "../tool-callback-config.js";
import type { BridgeToolCallbackConfig } from "../tool-callback-config.js";
import type { CliArgs } from "../types.js";

const READY_LINE_PREFIX = "cursor-sdk-bridge ready ";

async function main(): Promise<void> {
    const args = parseArgs(process.argv.slice(2));
    if (args.help) {
        process.stdout.write(helpText());
        return;
    }
    const port = parsePort(args.port ?? process.env.CURSOR_SDK_BRIDGE_PORT, "port");
    const host = args.host ?? process.env.CURSOR_SDK_BRIDGE_HOST ?? "127.0.0.1";
    const workspaceRef = args.workspace
        ? resolve(args.workspace)
        : process.env.CURSOR_SDK_BRIDGE_WORKSPACE;
    const stateRoot = args.stateRoot
        ? resolve(args.stateRoot)
        : process.env.CURSOR_SDK_BRIDGE_STATE_ROOT;
    setBridgeToolCallbackConfig(
        readToolCallbackConfigFromArgs(args) ?? readBridgeToolCallbackConfigFromEnv(),
    );
    configureDefaultLocalStore(args.localStore ?? process.env.CURSOR_SDK_LOCAL_STORE);
    const server = await startCursorSdkBridgeServer({
        port,
        host,
        workspaceRef,
        stateRoot,
        maxConcurrentAgents: args.maxConcurrentAgents,
        maxMessageBytes: args.maxMessageBytes,
    });
    const authTokenFile = writeAuthTokenFile(server.address.authToken);
    const readyAddress = {
        schemaVersion: server.address.schemaVersion,
        serverVersion: server.address.serverVersion,
        pid: server.address.pid,
        transport: server.address.transport,
        protocol: server.address.protocol,
        host: server.address.host,
        port: server.address.port,
        url: server.address.url,
        authTokenFile,
        workspaceRef: server.address.workspaceRef,
        stateRoot: server.address.stateRoot,
        maxConcurrentAgents: server.address.maxConcurrentAgents,
        maxMessageBytes: server.address.maxMessageBytes,
    };
    process.stderr.write(`${READY_LINE_PREFIX}${JSON.stringify(readyAddress)}\n`);
    const shutdown = async (signal) => {
        try {
            await server.close();
        }
        catch (err) {
            process.stderr.write(`cursor-sdk-bridge ${signal} shutdown error: ${err instanceof Error ? err.stack || err.message : String(err)}\n`);
            process.exit(1);
        }
        process.exit(0);
    };
    process.once("SIGINT", () => {
        void shutdown("SIGINT");
    });
    process.once("SIGTERM", () => {
        void shutdown("SIGTERM");
    });
}
main().catch(error => {
    process.stderr.write(`cursor-sdk-bridge failed: ${error instanceof Error ? error.stack || error.message : String(error)}\n`);
    process.exit(1);
});
function parseArgs(args: string[]): CliArgs {
    const parsed: CliArgs = {};
    for (let index = 0; index < args.length; index += 1) {
        const arg = args[index];
        switch (arg) {
            case "--host":
                parsed.host = takeValue(args, ++index, arg);
                break;
            case "--port":
                parsed.port = takeValue(args, ++index, arg);
                break;
            case "--workspace":
                parsed.workspace = takeValue(args, ++index, arg);
                break;
            case "--state-root":
                parsed.stateRoot = takeValue(args, ++index, arg);
                break;
            case "--tool-callback-url":
                parsed.toolCallbackUrl = takeValue(args, ++index, arg);
                break;
            case "--tool-callback-auth-token":
                parsed.toolCallbackAuthToken = takeValue(args, ++index, arg);
                break;
            case "--local-store":
                parsed.localStore = takeValue(args, ++index, arg);
                break;
            case "--max-concurrent-agents":
                parsed.maxConcurrentAgents = parsePositiveInteger({
                    value: takeValue(args, ++index, arg),
                    name: arg,
                });
                break;
            case "--max-message-bytes":
                parsed.maxMessageBytes = parsePositiveInteger({
                    value: takeValue(args, ++index, arg),
                    name: arg,
                });
                break;
            case "--help":
            case "-h":
                parsed.help = true;
                break;
            default:
                throw new Error(`Unknown argument: ${arg}`);
        }
    }
    return parsed;
}
function takeValue(args: string[], index: number, flag: string): string {
    const value = args[index];
    if (value === undefined || value.startsWith("-")) {
        throw new Error(`Missing value for ${flag}`);
    }
    return value;
}
function parsePort(value: string | undefined, name: string): number {
    if (value === undefined || value.length === 0) {
        return 0;
    }
    const port = Number.parseInt(value, 10);
    if (!Number.isInteger(port) || port < 0 || port > 65535) {
        throw new Error(`Invalid ${name}: ${value}`);
    }
    return port;
}
function parsePositiveInteger(options: { value: string; name: string }): number {
    const { value, name } = options;
    const parsed = Number.parseInt(value, 10);
    if (!Number.isInteger(parsed) || parsed <= 0) {
        throw new Error(`Invalid ${name}: ${value}`);
    }
    return parsed;
}
function readToolCallbackConfigFromArgs(args: CliArgs): BridgeToolCallbackConfig | undefined {
    const url = args.toolCallbackUrl?.trim() ?? process.env.CURSOR_SDK_TOOL_CALLBACK_URL?.trim();
    const authToken =
        args.toolCallbackAuthToken?.trim() ??
        process.env.CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN?.trim() ??
        process.env.CURSOR_SDK_TOOL_CALLBACK_TOKEN?.trim();
    if (!url && !authToken) {
        return undefined;
    }
    if (!url || !authToken) {
        throw new Error(
            "Tool callback configuration requires both URL and auth token " +
            "(--tool-callback-url / --tool-callback-auth-token or " +
            "CURSOR_SDK_TOOL_CALLBACK_URL / CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN)",
        );
    }
    return { url, authToken };
}

function configureDefaultLocalStore(raw: string | undefined): void {
    const trimmed = raw?.trim();
    if (!trimmed) {
        return;
    }
    let config;
    try {
        config = LocalAgentStoreConfig.fromJsonString(trimmed, {
            ignoreUnknownFields: true,
        });
    }
    catch (error) {
        throw new Error(`Invalid --local-store JSON: ${error instanceof Error ? error.message : String(error)}`);
    }
    const store = createLocalAgentStoreFromProtoConfig(config);
    if (store) {
        configureCursorSdk({ local: { store } });
    }
}
function writeAuthTokenFile(authToken: string): string {
    const dir = mkdtempSync(join(tmpdir(), "cursor-sdk-bridge-"));
    const path = join(dir, "auth-token");
    writeFileSync(path, authToken, { mode: 0o600 });
    return path;
}
function helpText() {
    return `Usage: cursor-sdk-bridge [options]

Options:
  --host <host>                         Host to bind (default: 127.0.0.1)
  --port <port>                         Port to bind (default: 0)
  --workspace <path>                    Workspace root used in discovery
  --state-root <path>                   Bridge state root used in discovery
  --tool-callback-url <url>             Host custom-tool callback Connect base URL
  --tool-callback-auth-token <token>    Bearer token for custom-tool callback RPCs
  --local-store <json>                  Default local.store config (LocalAgentStoreConfig
                                        JSON) used for every agent operation
  --max-concurrent-agents <count>       Advertised agent concurrency limit
  --max-message-bytes <bytes>           Advertised max message size
  --help, -h                            Show this help

Environment variables:
  CURSOR_SDK_BRIDGE_PORT                Same as --port
  CURSOR_SDK_BRIDGE_HOST                Same as --host
  CURSOR_SDK_BRIDGE_WORKSPACE           Same as --workspace
  CURSOR_SDK_BRIDGE_STATE_ROOT          Same as --state-root
  CURSOR_SDK_TOOL_CALLBACK_URL          Same as --tool-callback-url
  CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN   Same as --tool-callback-auth-token
  CURSOR_SDK_LOCAL_STORE                Same as --local-store
`;
}
