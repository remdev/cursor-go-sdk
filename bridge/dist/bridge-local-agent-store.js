import { Code, ConnectError } from "@connectrpc/connect";
import { JsonlLocalAgentStore } from "@cursor/sdk";

/**
 * Build a default local agent store from proto LocalAgentStoreConfig.
 * Go SDK supports jsonl and default sqlite only — not Python host store callbacks.
 */
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
            throw new ConnectError('local.store type "custom" (host store callback) is not supported in cursor-go-sdk bridge', Code.Unimplemented);
        default:
            throw new ConnectError(`Unsupported local.store type: ${store.type}`, Code.InvalidArgument);
    }
}
