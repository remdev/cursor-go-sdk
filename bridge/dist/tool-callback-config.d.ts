export interface BridgeToolCallbackConfig {
    readonly url: string;
    readonly authToken: string;
}
export declare function setBridgeToolCallbackConfig(config: BridgeToolCallbackConfig | undefined): void;
export declare function getBridgeToolCallbackConfig(): BridgeToolCallbackConfig | undefined;
export declare function readBridgeToolCallbackConfigFromEnv(): BridgeToolCallbackConfig | undefined;
//# sourceMappingURL=tool-callback-config.d.ts.map