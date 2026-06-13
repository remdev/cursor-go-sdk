import { SdkStoreCallbackService } from "@anysphere/proto/sdk/v1/sdk_store_callback_service_connect.js";
import { CallStoreRequest } from "@anysphere/proto/sdk/v1/sdk_store_callback_service_pb.js";
import { Struct } from "@bufbuild/protobuf";
import { Code, ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { composeLocalAgentStore, JsonlLocalAgentStore, } from "@cursor/sdk";
import { getBridgeStoreCallbackConfig } from "./store-callback-config.js";
let callbackClient;
function getStoreCallbackClient() {
    const config = getBridgeStoreCallbackConfig();
    if (!config) {
        throw new ConnectError("Store callbacks are not configured on this bridge (set CURSOR_SDK_STORE_CALLBACK_URL and CURSOR_SDK_STORE_CALLBACK_AUTH_TOKEN, or pass --store-callback-url / --store-callback-auth-token)", Code.FailedPrecondition);
    }
    if (!callbackClient) {
        const transport = createConnectTransport({
            baseUrl: config.url,
            httpVersion: "1.1",
            // The host server (e.g. the Python SDK) speaks Connect/JSON, so the
            // request/response bodies must be JSON rather than binary protobuf.
            useBinaryFormat: false,
            interceptors: [
                next => async request => {
                    request.header.set("Authorization", `Bearer ${config.authToken}`);
                    return await next(request);
                },
            ],
        });
        callbackClient = createClient(SdkStoreCallbackService, transport);
    }
    return callbackClient;
}
/**
 * Forward a single LocalAgentStore operation to a host-owned substore over the
 * loopback `CallStore` RPC. The input/output objects ride as `Struct`; an unset
 * output is treated as `null`.
 */
async function callHostStore(substore, method, input) {
    const response = await getStoreCallbackClient().callStore(new CallStoreRequest({
        substore,
        method,
        // JSON.stringify drops `undefined` fields, keeping Struct.fromJson happy.
        input: Struct.fromJsonString(JSON.stringify(input !== null && input !== void 0 ? input : {})),
    }));
    return (response.output ? response.output.toJson() : null);
}
function bytesToBase64(data) {
    return Buffer.from(data).toString("base64");
}
function base64ToBytes(value) {
    return new Uint8Array(Buffer.from(value, "base64"));
}
function createHostAgents() {
    return {
        get: input => callHostStore("agents", "get", input),
        create: input => callHostStore("agents", "create", input),
        update: input => callHostStore("agents", "update", input),
        delete: async input => {
            await callHostStore("agents", "delete", input);
        },
        list: input => callHostStore("agents", "list", input !== null && input !== void 0 ? input : {}),
    };
}
function createHostRuns() {
    return {
        get: input => callHostStore("runs", "get", input),
        create: input => callHostStore("runs", "create", input),
        update: input => callHostStore("runs", "update", input),
        delete: async input => {
            await callHostStore("runs", "delete", input);
        },
        list: input => callHostStore("runs", "list", input !== null && input !== void 0 ? input : {}),
    };
}
function createHostRunEvents() {
    return {
        append: input => callHostStore("runEvents", "append", input),
        list: input => callHostStore("runEvents", "list", input),
        delete: async input => {
            await callHostStore("runEvents", "delete", input);
        },
    };
}
function createHostCheckpoints() {
    return {
        get: async input => {
            const response = await callHostStore("checkpoints", "get", {
                agentId: input.agentId,
                blobId: input.blobId,
            });
            if (!(response === null || response === void 0 ? void 0 : response.found) || typeof response.data !== "string") {
                return null;
            }
            return base64ToBytes(response.data);
        },
        create: async input => {
            await callHostStore("checkpoints", "create", {
                agentId: input.agentId,
                blobId: input.blobId,
                data: bytesToBase64(input.data),
            });
        },
        update: async input => {
            await callHostStore("checkpoints", "update", {
                agentId: input.agentId,
                blobId: input.blobId,
                data: bytesToBase64(input.data),
            });
        },
        delete: async input => {
            await callHostStore("checkpoints", "delete", input);
        },
        list: input => callHostStore("checkpoints", "list", input !== null && input !== void 0 ? input : {}),
    };
}
/**
 * A {@link LocalAgentStore} whose every substore is implemented by the host
 * process (e.g. the Python SDK) over the loopback `CallStore` RPC. Used for
 * `local.store` configs of type `custom`.
 */
export function createCallbackLocalAgentStore() {
    return composeLocalAgentStore({
        agents: createHostAgents(),
        runs: createHostRuns(),
        runEvents: createHostRunEvents(),
        checkpoints: createHostCheckpoints(),
    });
}
export function createLocalAgentStoreFromProtoConfig(store) {
    var _a;
    if (!(store === null || store === void 0 ? void 0 : store.type)) {
        return undefined;
    }
    switch (store.type) {
        case "sqlite":
            // Default bridge-owned SQLite store (same as leaving local.store unset).
            return undefined;
        case "jsonl": {
            const rootDir = (_a = store.rootDir) === null || _a === void 0 ? void 0 : _a.trim();
            if (!rootDir) {
                throw new ConnectError('local.store type "jsonl" requires rootDir', Code.InvalidArgument);
            }
            return new JsonlLocalAgentStore(rootDir);
        }
        case "custom":
            // Fully host-owned store: every substore (agents, runs, runEvents,
            // checkpoints) is implemented by the host process over loopback RPC.
            return createCallbackLocalAgentStore();
        default:
            throw new ConnectError(`Unsupported local.store type: ${store.type}`, Code.InvalidArgument);
    }
}
export function resetBridgeStoreCallbackClientForTests() {
    callbackClient = undefined;
}
