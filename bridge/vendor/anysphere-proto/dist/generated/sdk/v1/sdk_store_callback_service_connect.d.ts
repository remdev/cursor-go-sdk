import { CallStoreRequest, CallStoreResponse } from "./sdk_store_callback_service_pb.js";
import { MethodKind } from "@bufbuild/protobuf";
/**
 * Store callbacks from the Node bridge to a Python (or other) host process.
 *
 * The bridge forwards local-agent-store operations to the host over a single
 * generic `CallStore` RPC, mirroring the `@cursor/sdk` `LocalAgentStore`
 * topology (substores: agents, runs, runEvents, checkpoints). Inputs/outputs
 * are structured JSON so the bridge can forward the full surface without a
 * typed message per operation. This powers `local.store` type "custom", where
 * every substore is implemented by the host process.
 *
 * @generated from service sdk.v1.SdkStoreCallbackService
 */
export declare const SdkStoreCallbackService: {
    readonly typeName: "sdk.v1.SdkStoreCallbackService";
    readonly methods: {
        /**
         * @generated from rpc sdk.v1.SdkStoreCallbackService.CallStore
         */
        readonly callStore: {
            readonly name: "CallStore";
            readonly I: typeof CallStoreRequest;
            readonly O: typeof CallStoreResponse;
            readonly kind: MethodKind.Unary;
        };
    };
};
//# sourceMappingURL=sdk_store_callback_service_connect.d.ts.map