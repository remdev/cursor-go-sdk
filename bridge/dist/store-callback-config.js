let storeCallbackConfig;
export function setBridgeStoreCallbackConfig(config) {
    storeCallbackConfig = config;
}
export function getBridgeStoreCallbackConfig() {
    return storeCallbackConfig;
}
export function readBridgeStoreCallbackConfigFromEnv() {
    var _a, _b, _c;
    var _d;
    const url = (_a = process.env.CURSOR_SDK_STORE_CALLBACK_URL) === null || _a === void 0 ? void 0 : _a.trim();
    const authToken = (_d = (_b = process.env.CURSOR_SDK_STORE_CALLBACK_AUTH_TOKEN) === null || _b === void 0 ? void 0 : _b.trim()) !== null && _d !== void 0 ? _d : (_c = process.env.CURSOR_SDK_STORE_CALLBACK_TOKEN) === null || _c === void 0 ? void 0 : _c.trim();
    if (!url || !authToken) {
        return undefined;
    }
    return { url, authToken };
}
