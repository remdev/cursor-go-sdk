export interface BridgeStoreCallbackConfig {
    readonly url: string;
    readonly authToken: string;
}
export declare function setBridgeStoreCallbackConfig(config: BridgeStoreCallbackConfig | undefined): void;
export declare function getBridgeStoreCallbackConfig(): BridgeStoreCallbackConfig | undefined;
export declare function readBridgeStoreCallbackConfigFromEnv(): BridgeStoreCallbackConfig | undefined;
//# sourceMappingURL=store-callback-config.d.ts.map