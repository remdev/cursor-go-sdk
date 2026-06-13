export interface BridgeToolCallbackConfig {
    url: string;
    authToken: string;
}

let toolCallbackConfig: BridgeToolCallbackConfig | undefined;

export function setBridgeToolCallbackConfig(config: BridgeToolCallbackConfig | undefined): void {
    toolCallbackConfig = config;
}

export function getBridgeToolCallbackConfig(): BridgeToolCallbackConfig | undefined {
    return toolCallbackConfig;
}

export function readBridgeToolCallbackConfigFromEnv(): BridgeToolCallbackConfig | undefined {
    const url = process.env.CURSOR_SDK_TOOL_CALLBACK_URL?.trim();
    const authToken =
        process.env.CURSOR_SDK_TOOL_CALLBACK_AUTH_TOKEN?.trim() ??
        process.env.CURSOR_SDK_TOOL_CALLBACK_TOKEN?.trim();
    if (!url || !authToken) {
        return undefined;
    }
    return { url, authToken };
}
