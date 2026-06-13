import type { CustomToolDefinition } from "@anysphere/proto/sdk/v1/sdk_messages_pb.js";
import type { SDKCustomTool, SDKJsonValue, LocalAgentOptions as SdkLocalAgentOptions } from "@cursor/sdk";
export type HostCustomToolDefinition = {
    description?: string;
    inputSchema?: Record<string, SDKJsonValue>;
};
export declare function hasCustomToolDefinitions(customTools: Record<string, CustomToolDefinition> | undefined): boolean;
export declare function protoCustomToolDefinitionsToHost(customTools: Record<string, CustomToolDefinition> | undefined): Record<string, HostCustomToolDefinition> | undefined;
export declare function buildHostBackedCustomTools(agentId: string, definitions: Record<string, HostCustomToolDefinition>): Record<string, SDKCustomTool>;
export declare function attachHostCustomToolsToSdkLocalOptions(local: SdkLocalAgentOptions | undefined, definitions: Record<string, HostCustomToolDefinition> | undefined, agentId: string | undefined): SdkLocalAgentOptions | undefined;
export declare function assertCustomToolsConfigured(customTools: Record<string, CustomToolDefinition> | undefined): void;
export declare function assertCustomToolsLocalOnly(options: {
    cloud?: unknown;
    localCustomTools?: Record<string, CustomToolDefinition>;
}): void;
export declare function resetBridgeToolCallbackClient(): void;
export declare function resetBridgeToolCallbackClientForTests(): void;
//# sourceMappingURL=bridge-custom-tools.d.ts.map