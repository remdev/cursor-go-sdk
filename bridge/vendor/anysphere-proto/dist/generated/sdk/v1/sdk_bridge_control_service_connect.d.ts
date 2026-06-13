import { GetVersionRequest, GetVersionResponse, PingRequest, PingResponse, SetToolCallbackRequest, SetToolCallbackResponse, ShutdownRequest, ShutdownResponse } from "./sdk_bridge_control_service_pb.js";
import { MethodKind } from "@bufbuild/protobuf";
/**
 * Bridge-local health, version, and lifecycle APIs.
 *
 * @generated from service sdk.v1.SdkBridgeControlService
 */
export declare const SdkBridgeControlService: {
    readonly typeName: "sdk.v1.SdkBridgeControlService";
    readonly methods: {
        /**
         * @generated from rpc sdk.v1.SdkBridgeControlService.Ping
         */
        readonly ping: {
            readonly name: "Ping";
            readonly I: typeof PingRequest;
            readonly O: typeof PingResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkBridgeControlService.Shutdown
         */
        readonly shutdown: {
            readonly name: "Shutdown";
            readonly I: typeof ShutdownRequest;
            readonly O: typeof ShutdownResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkBridgeControlService.GetVersion
         */
        readonly getVersion: {
            readonly name: "GetVersion";
            readonly I: typeof GetVersionRequest;
            readonly O: typeof GetVersionResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * Point the bridge at a host-process custom-tool callback server after
         * startup. Lets SDK clients that attach to an already-running bridge
         * (Client.connect) register host execute handlers. Same-host (loopback) only.
         *
         * @generated from rpc sdk.v1.SdkBridgeControlService.SetToolCallback
         */
        readonly setToolCallback: {
            readonly name: "SetToolCallback";
            readonly I: typeof SetToolCallbackRequest;
            readonly O: typeof SetToolCallbackResponse;
            readonly kind: MethodKind.Unary;
        };
    };
};
//# sourceMappingURL=sdk_bridge_control_service_connect.d.ts.map