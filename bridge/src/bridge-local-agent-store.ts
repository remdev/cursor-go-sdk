import { Code, ConnectError } from "@connectrpc/connect";
import { JsonlLocalAgentStore } from "@cursor/sdk";
import type { LocalAgentStoreConfig } from "../gen/ts/sdk/v1/sdk_messages_pb.js";

/**
 * Build a default local agent store from proto LocalAgentStoreConfig.
 * Go SDK supports jsonl and default sqlite only — host store callbacks over RPC are not supported.
 */
export function createLocalAgentStoreFromProtoConfig(
    store: LocalAgentStoreConfig | undefined,
): JsonlLocalAgentStore | undefined {
    if (!store?.type) {
        return undefined;
    }
    switch (store.type) {
        case "sqlite":
            // Default bridge-owned SQLite store (same as leaving local.store unset).
            return undefined;
        case "jsonl": {
            const rootDir = store.rootDir?.trim();
            if (!rootDir) {
                throw new ConnectError(
                    'local.store type "jsonl" requires rootDir',
                    Code.InvalidArgument,
                );
            }
            return new JsonlLocalAgentStore(rootDir);
        }
        case "custom":
            throw new ConnectError(
                'local.store type "custom" (host store callback) is not supported in cursor-go-sdk bridge',
                Code.Unimplemented,
            );
        default:
            throw new ConnectError(
                `Unsupported local.store type: ${store.type}`,
                Code.InvalidArgument,
            );
    }
}
