import type { LocalAgentStoreConfig } from "@anysphere/proto/sdk/v1/sdk_messages_pb.js";
import { type LocalAgentStore } from "@cursor/sdk";
/**
 * A {@link LocalAgentStore} whose every substore is implemented by the host
 * process (e.g. the Python SDK) over the loopback `CallStore` RPC. Used for
 * `local.store` configs of type `custom`.
 */
export declare function createCallbackLocalAgentStore(): LocalAgentStore;
export declare function createLocalAgentStoreFromProtoConfig(store?: LocalAgentStoreConfig): LocalAgentStore | undefined;
export declare function resetBridgeStoreCallbackClientForTests(): void;
//# sourceMappingURL=bridge-local-agent-store.d.ts.map