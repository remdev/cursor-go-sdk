#!/usr/bin/env node
import { mkdtempSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { LocalAgentStoreConfig } from "@anysphere/proto/sdk/v1/sdk_messages_pb.js";
import { configureCursorSdk } from "@cursor/sdk";
import { createLocalAgentStoreFromProtoConfig } from "../bridge-local-agent-store.js";
import { installProcessErrorSurvivors, SURVIVE_UNCAUGHT_ENV_VAR, shouldInstallProcessErrorSurvivors, } from "../process-error-survivors.js";
import { startCursorSdkBridgeServer } from "../server.js";
import { readBridgeStoreCallbackConfigFromEnv, setBridgeStoreCallbackConfig, } from "../store-callback-config.js";
import { readBridgeToolCallbackConfigFromEnv, setBridgeToolCallbackConfig, } from "../tool-callback-config.js";
const READY_LINE_PREFIX = "cursor-sdk-bridge ready ";
async function main() {
    var _a, _b, _c, _d, _e, _f;
    // Default behavior is Node's default: uncaughtException kills the process,
    // surfacing the leak loudly. Only the e2e suite (or any other consciously
    // opt-in caller) sets `CURSOR_SDK_BRIDGE_SURVIVE_UNCAUGHT=1` to trade
    // visibility for liveness. See `process-error-survivors.ts` for the full
    // rationale and the removal-criteria TODO.
    if (shouldInstallProcessErrorSurvivors(process.env)) {
        installProcessErrorSurvivors();
        process.stderr.write(`cursor-sdk-bridge: ${SURVIVE_UNCAUGHT_ENV_VAR}=1 — installing process-level error survivors (band-aid mode)\n`);
    }
    const args = parseArgs(process.argv.slice(2));
    if (args.help) {
        process.stdout.write(helpText());
        return;
    }
    const port = parsePort((_a = args.port) !== null && _a !== void 0 ? _a : process.env.CURSOR_SDK_BRIDGE_PORT, "port");
    const host = (_c = (_b = args.host) !== null && _b !== void 0 ? _b : process.env.CURSOR_SDK_BRIDGE_HOST) !== null && _c !== void 0 ? _c : "127.0.0.1";
    const workspaceRef = args.workspace
        ? resolve(args.workspace)
        : process.env.CURSOR_SDK_BRIDGE_WORKSPACE;
    const stateRoot = args.stateRoot
        ? resolve(args.stateRoot)
        : process.env.CURSOR_SDK_BRIDGE_STATE_ROOT;
    setBridgeStoreCallbackConfig((_d = readStoreCallbackConfigFromArgs(args)) !== null && _d !== void 0 ? _d : readBridgeStoreCallbackConfigFromEnv());
    setBridgeToolCallbackConfig((_e = readToolCallbackConfigFromArgs(args)) !== null && _e !== void 0 ? _e : readBridgeToolCallbackConfigFromEnv());
    // Bind the client-configured local store once so every agent operation
    // (create / list / get / resume) uses it, mirroring the TS SDK's
    // `configureCursorSdk({ local: { store } })`. Per-call `local.store` still
    // overrides on create/resume.
    configureDefaultLocalStore((_f = args.localStore) !== null && _f !== void 0 ? _f : process.env.CURSOR_SDK_LOCAL_STORE);
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
    // Shutdown must catch its own errors and force exit, otherwise a thrown
    // `server.close()` (e.g., `registry.dispose()` failing during agent
    // teardown) would surface as an unhandled rejection. With the survivors
    // installed (opt-in, `CURSOR_SDK_BRIDGE_SURVIVE_UNCAUGHT=1`), that
    // rejection would be swallowed and the bridge would be unkillable via
    // SIGTERM/SIGINT until the SDK falls back to SIGKILL after a 5s timeout.
    // Catching here keeps the signal path robust regardless of whether the
    // survivors are installed.
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
function parseArgs(args) {
    const parsed = {};
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
            case "--store-callback-url":
                parsed.storeCallbackUrl = takeValue(args, ++index, arg);
                break;
            case "--store-callback-auth-token":
                parsed.storeCallbackAuthToken = takeValue(args, ++index, arg);
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
function takeValue(args, index, flag) {
    const value = args[index];
    if (value === undefined || value.startsWith("-")) {
        throw new Error(`Missing value for ${flag}`);
    }
    return value;
}
function parsePort(value, name) {
    if (value === undefined || value.length === 0) {
        return 0;
    }
    const port = Number.parseInt(value, 10);
    if (!Number.isInteger(port) || port < 0 || port > 65535) {
        throw new Error(`Invalid ${name}: ${value}`);
    }
    return port;
}
function parsePositiveInteger(options) {
    const { value, name } = options;
    const parsed = Number.parseInt(value, 10);
    if (!Number.isInteger(parsed) || parsed <= 0) {
        throw new Error(`Invalid ${name}: ${value}`);
    }
    return parsed;
}
function readStoreCallbackConfigFromArgs(args) {
    var _a, _b, _c, _d, _e;
    var _f, _g, _h;
    const url = (_f = (_a = args.storeCallbackUrl) === null || _a === void 0 ? void 0 : _a.trim()) !== null && _f !== void 0 ? _f : (_b = process.env.CURSOR_SDK_STORE_CALLBACK_URL) === null || _b === void 0 ? void 0 : _b.trim();
    const authToken = (_h = (_g = (_c = args.storeCallbackAuthToken) === null || _c === void 0 ? void 0 : _c.trim()) !== null && _g !== void 0 ? _g : (_d = process.env.CURSOR_SDK_STORE_CALLBACK_AUTH_TOKEN) === null || _d === void 0 ? void 0 : _d.trim()) !== null && _h !== void 0 ? _h : (_e = process.env.CURSOR_SDK_STORE_CALLBACK_TOKEN) === null || _e === void 0 ? void 0 : _e.trim();
    if (!url && !authToken) {
        return undefined;
    }
    if (!url || !authToken) {
        throw new Error("Store callback configuration requires both URL and auth token (--store-callback-url / --store-callback-auth-token or CURSOR_SDK_STORE_CALLBACK_URL / CURSOR_SDK_STORE_CALLBACK_AUTH_TOKEN)");
    }
    return { url, authToken };
}
function readToolCallbackConfigFromArgs(args) {
    var _a, _b, _c, _d, _e;
    var _f, _g, _h;
    const url = (_f = (_a = args.toolCallbackUrl) === null || _a === void 0 ? void 0 : _a.trim()) !== null && _f !== void 0 ? _f : (_b = process.env.CURSOR_SDK_TOOL_CALLBACK_URL) === null || _b === void 0 ? void 0 : _b.trim();
    const authToken = (_h = (_g = (_c = args.toolCallbackAuthToken) === null || _c === void 0 ? void 0 : _c.trim()) !== null && _g !== void 0 ? _g : (_d = process.env.CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN) === null || _d === void 0 ? void 0 : _d.trim()) !== null && _h !== void 0 ? _h : (_e = process.env.CURSOR_SDK_TOOL_CALLBACK_TOKEN) === null || _e === void 0 ? void 0 : _e.trim();
    if (!url && !authToken) {
        return undefined;
    }
    if (!url || !authToken) {
        throw new Error("Tool callback configuration requires both URL and auth token (--tool-callback-url / --tool-callback-auth-token or CURSOR_SDK_TOOL_CALLBACK_URL / CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN)");
    }
    return { url, authToken };
}
function configureDefaultLocalStore(raw) {
    const trimmed = raw === null || raw === void 0 ? void 0 : raw.trim();
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
function writeAuthTokenFile(authToken) {
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
  --store-callback-url <url>            Python/host store callback Connect base URL
  --store-callback-auth-token <token>   Bearer token for store callback RPCs
  --tool-callback-url <url>             Python/host custom-tool callback Connect base URL
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
  CURSOR_SDK_STORE_CALLBACK_URL         Same as --store-callback-url
  CURSOR_SDK_STORE_CALLBACK_AUTH_TOKEN  Same as --store-callback-auth-token
  CURSOR_SDK_TOOL_CALLBACK_URL          Same as --tool-callback-url
  CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN   Same as --tool-callback-auth-token
  CURSOR_SDK_LOCAL_STORE                Same as --local-store
  ${SURVIVE_UNCAUGHT_ENV_VAR}=1
                                        Opt-in: install process-level
                                        \`uncaughtException\` and
                                        \`unhandledRejection\` survivors that
                                        log and continue. OFF by default
                                        (Node's standard fail-loud behavior)
                                        because surviving may leave the
                                        in-memory registry / IPC channels in
                                        an undefined state. Only enable for
                                        e2e suites that need the bridge to
                                        outlive a leaked handler.
`;
}
