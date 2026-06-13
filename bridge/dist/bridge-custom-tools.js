import { SdkCustomToolCallbackService } from "@anysphere/proto/sdk/v1/sdk_custom_tool_callback_service_connect.js";
import { CallCustomToolRequest } from "@anysphere/proto/sdk/v1/sdk_custom_tool_callback_service_pb.js";
import { Struct } from "@bufbuild/protobuf";
import { Code, ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { getBridgeToolCallbackConfig } from "./tool-callback-config.js";
let callbackClient;
function getToolCallbackClient() {
    const config = getBridgeToolCallbackConfig();
    if (!config) {
        throw new ConnectError("Custom tool callbacks are not configured on this bridge (pass --tool-callback-url / --tool-callback-auth-token)", Code.FailedPrecondition);
    }
    if (!callbackClient) {
        const transport = createConnectTransport({
            baseUrl: config.url,
            httpVersion: "1.1",
            useBinaryFormat: false,
            interceptors: [
                next => async request => {
                    request.header.set("Authorization", `Bearer ${config.authToken}`);
                    return await next(request);
                },
            ],
        });
        callbackClient = createClient(SdkCustomToolCallbackService, transport);
    }
    return callbackClient;
}
export function hasCustomToolDefinitions(customTools) {
    return customTools !== undefined && Object.keys(customTools).length > 0;
}
export function protoCustomToolDefinitionsToHost(customTools) {
    if (!hasCustomToolDefinitions(customTools)) {
        return undefined;
    }
    const hostTools = {};
    for (const [toolName, definition] of Object.entries(customTools)) {
        hostTools[toolName] = {
            description: definition.description || undefined,
            inputSchema: definition.inputSchema
                ? definition.inputSchema.toJson()
                : undefined,
        };
    }
    return hostTools;
}
async function callHostCustomTool(agentId, toolName, args, context) {
    const response = await getToolCallbackClient().callCustomTool(new CallCustomToolRequest({
        agentId,
        toolName,
        args: Struct.fromJson(args),
        toolCallId: context.toolCallId,
    }));
    if (!response.result) {
        throw new Error(`Host custom tool ${toolName} returned no result`);
    }
    return response.result.toJson();
}
export function buildHostBackedCustomTools(agentId, definitions) {
    const customTools = {};
    for (const [toolName, definition] of Object.entries(definitions)) {
        customTools[toolName] = {
            description: definition.description,
            inputSchema: definition.inputSchema,
            execute: async (args, context) => callHostCustomTool(agentId, toolName, args, context),
        };
    }
    return customTools;
}
export function attachHostCustomToolsToSdkLocalOptions(local, definitions, agentId) {
    if (!definitions || Object.keys(definitions).length === 0) {
        return local;
    }
    if (!agentId) {
        throw new ConnectError("agent_id is required to attach host-backed custom tools", Code.Internal);
    }
    return Object.assign(Object.assign({}, (local !== null && local !== void 0 ? local : {})), { customTools: buildHostBackedCustomTools(agentId, definitions) });
}
export function assertCustomToolsConfigured(customTools) {
    if (!hasCustomToolDefinitions(customTools)) {
        return;
    }
    if (!getBridgeToolCallbackConfig()) {
        throw new ConnectError("local.custom_tools requires a host tool callback server (--tool-callback-url / --tool-callback-auth-token)", Code.FailedPrecondition);
    }
}
export function assertCustomToolsLocalOnly(options) {
    if (!options.cloud) {
        return;
    }
    if (hasCustomToolDefinitions(options.localCustomTools)) {
        throw new ConnectError("local.custom_tools is only supported for local SDK agents", Code.InvalidArgument);
    }
}
// Drop the cached callback client so the next call rebuilds it against the
// current BridgeToolCallbackConfig. Call after setBridgeToolCallbackConfig
// changes the endpoint at runtime (e.g. SetToolCallback from Client.connect).
export function resetBridgeToolCallbackClient() {
    callbackClient = undefined;
}
export function resetBridgeToolCallbackClientForTests() {
    resetBridgeToolCallbackClient();
}
