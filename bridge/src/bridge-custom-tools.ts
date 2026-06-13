import { SdkCustomToolCallbackService } from "../gen/ts/sdk/v1/sdk_custom_tool_callback_service_connect.js";
import { CallCustomToolRequest } from "../gen/ts/sdk/v1/sdk_custom_tool_callback_service_pb.js";
import { CustomToolDefinition } from "../gen/ts/sdk/v1/sdk_messages_pb.js";
import { Struct, type JsonValue } from "@bufbuild/protobuf";
import {
    Code,
    ConnectError,
    createClient,
    type Client,
} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { getBridgeToolCallbackConfig } from "./tool-callback-config.js";

type CustomToolCallbackClient = Client<typeof SdkCustomToolCallbackService>;

let callbackClient: CustomToolCallbackClient | undefined;

function getToolCallbackClient(): CustomToolCallbackClient {
    const config = getBridgeToolCallbackConfig();
    if (!config) {
        throw new ConnectError(
            "Custom tool callbacks are not configured on this bridge (pass --tool-callback-url / --tool-callback-auth-token)",
            Code.FailedPrecondition,
        );
    }
    if (!callbackClient) {
        const transport = createConnectTransport({
            baseUrl: config.url,
            httpVersion: "1.1",
            useBinaryFormat: false,
            interceptors: [
                (next) => async (request) => {
                    request.header.set("Authorization", `Bearer ${config.authToken}`);
                    return await next(request);
                },
            ],
        });
        callbackClient = createClient(SdkCustomToolCallbackService, transport);
    }
    return callbackClient;
}

export function hasCustomToolDefinitions(
    customTools: { [key: string]: CustomToolDefinition } | undefined,
): boolean {
    return customTools !== undefined && Object.keys(customTools).length > 0;
}

export function protoCustomToolDefinitionsToHost(
    customTools: { [key: string]: CustomToolDefinition } | undefined,
) {
    if (!hasCustomToolDefinitions(customTools)) {
        return undefined;
    }
    const hostTools: Record<string, { description?: string; inputSchema?: unknown }> = {};
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

async function callHostCustomTool(
    agentId: string,
    toolName: string,
    args: JsonValue,
    context: { toolCallId?: string },
) {
    const response = await getToolCallbackClient().callCustomTool(
        new CallCustomToolRequest({
            agentId,
            toolName,
            args: Struct.fromJson(args),
            toolCallId: context.toolCallId,
        }),
    );
    if (!response.result) {
        throw new Error(`Host custom tool ${toolName} returned no result`);
    }
    return response.result.toJson();
}

export function buildHostBackedCustomTools(
    agentId: string,
    definitions: Record<string, { description?: string; inputSchema?: unknown }>,
) {
    const customTools: Record<
        string,
        {
            description?: string;
            inputSchema?: unknown;
            execute: (
                args: JsonValue,
                context: { toolCallId?: string },
            ) => Promise<unknown>;
        }
    > = {};
    for (const [toolName, definition] of Object.entries(definitions)) {
        customTools[toolName] = {
            description: definition.description,
            inputSchema: definition.inputSchema,
            execute: async (args, context) =>
                callHostCustomTool(agentId, toolName, args, context),
        };
    }
    return customTools;
}

export function attachHostCustomToolsToSdkLocalOptions<T extends Record<string, unknown>>(
    local: T | undefined,
    definitions: Record<string, { description?: string; inputSchema?: unknown }> | undefined,
    agentId: string | undefined,
) {
    if (!definitions || Object.keys(definitions).length === 0) {
        return local;
    }
    if (!agentId) {
        throw new ConnectError(
            "agent_id is required to attach host-backed custom tools",
            Code.Internal,
        );
    }
    return {
        ...(local ?? {}),
        customTools: buildHostBackedCustomTools(agentId, definitions),
    };
}

export function assertCustomToolsConfigured(
    customTools: { [key: string]: CustomToolDefinition } | undefined,
): void {
    if (!hasCustomToolDefinitions(customTools)) {
        return;
    }
    if (!getBridgeToolCallbackConfig()) {
        throw new ConnectError(
            "local.custom_tools requires a host tool callback server (--tool-callback-url / --tool-callback-auth-token)",
            Code.FailedPrecondition,
        );
    }
}

export function assertCustomToolsLocalOnly(options: {
    cloud?: boolean;
    localCustomTools?: { [key: string]: CustomToolDefinition };
}): void {
    if (!options.cloud) {
        return;
    }
    if (hasCustomToolDefinitions(options.localCustomTools)) {
        throw new ConnectError(
            "local.custom_tools is only supported for local SDK agents",
            Code.InvalidArgument,
        );
    }
}

// Drop the cached callback client so the next call rebuilds it against the
// current BridgeToolCallbackConfig. Call after setBridgeToolCallbackConfig
// changes the endpoint at runtime (e.g. SetToolCallback from Client.connect).
export function resetBridgeToolCallbackClient(): void {
    callbackClient = undefined;
}

export function resetBridgeToolCallbackClientForTests(): void {
    resetBridgeToolCallbackClient();
}
