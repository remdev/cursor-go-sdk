import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
/**
 * @generated from message sdk.v1.PingRequest
 */
export declare class PingRequest extends Message<PingRequest> {
    constructor(data?: PartialMessage<PingRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.PingRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PingRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PingRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PingRequest;
    static equals(a: PingRequest | PlainMessage<PingRequest> | undefined, b: PingRequest | PlainMessage<PingRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.PingResponse
 */
export declare class PingResponse extends Message<PingResponse> {
    /**
     * @generated from field: string message = 1;
     */
    message: string;
    constructor(data?: PartialMessage<PingResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.PingResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PingResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PingResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PingResponse;
    static equals(a: PingResponse | PlainMessage<PingResponse> | undefined, b: PingResponse | PlainMessage<PingResponse> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ShutdownRequest
 */
export declare class ShutdownRequest extends Message<ShutdownRequest> {
    /**
     * @generated from field: uint32 grace_seconds = 1;
     */
    graceSeconds: number;
    constructor(data?: PartialMessage<ShutdownRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ShutdownRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ShutdownRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ShutdownRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ShutdownRequest;
    static equals(a: ShutdownRequest | PlainMessage<ShutdownRequest> | undefined, b: ShutdownRequest | PlainMessage<ShutdownRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ShutdownResponse
 */
export declare class ShutdownResponse extends Message<ShutdownResponse> {
    constructor(data?: PartialMessage<ShutdownResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ShutdownResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ShutdownResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ShutdownResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ShutdownResponse;
    static equals(a: ShutdownResponse | PlainMessage<ShutdownResponse> | undefined, b: ShutdownResponse | PlainMessage<ShutdownResponse> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.GetVersionRequest
 */
export declare class GetVersionRequest extends Message<GetVersionRequest> {
    constructor(data?: PartialMessage<GetVersionRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.GetVersionRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetVersionRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetVersionRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetVersionRequest;
    static equals(a: GetVersionRequest | PlainMessage<GetVersionRequest> | undefined, b: GetVersionRequest | PlainMessage<GetVersionRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.GetVersionResponse
 */
export declare class GetVersionResponse extends Message<GetVersionResponse> {
    /**
     * @generated from field: string bridge_version = 1;
     */
    bridgeVersion: string;
    /**
     * @generated from field: string protocol_version = 2;
     */
    protocolVersion: string;
    /**
     * @generated from field: repeated string capabilities = 3;
     */
    capabilities: string[];
    constructor(data?: PartialMessage<GetVersionResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.GetVersionResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetVersionResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetVersionResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetVersionResponse;
    static equals(a: GetVersionResponse | PlainMessage<GetVersionResponse> | undefined, b: GetVersionResponse | PlainMessage<GetVersionResponse> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SetToolCallbackRequest
 */
export declare class SetToolCallbackRequest extends Message<SetToolCallbackRequest> {
    /**
     * Connect base URL of the host custom-tool callback server. Empty clears the
     * currently configured callback.
     *
     * @generated from field: string url = 1;
     */
    url: string;
    /**
     * Bearer token the bridge presents on CallCustomTool RPCs.
     *
     * @generated from field: string auth_token = 2;
     */
    authToken: string;
    constructor(data?: PartialMessage<SetToolCallbackRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SetToolCallbackRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SetToolCallbackRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SetToolCallbackRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SetToolCallbackRequest;
    static equals(a: SetToolCallbackRequest | PlainMessage<SetToolCallbackRequest> | undefined, b: SetToolCallbackRequest | PlainMessage<SetToolCallbackRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SetToolCallbackResponse
 */
export declare class SetToolCallbackResponse extends Message<SetToolCallbackResponse> {
    constructor(data?: PartialMessage<SetToolCallbackResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SetToolCallbackResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SetToolCallbackResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SetToolCallbackResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SetToolCallbackResponse;
    static equals(a: SetToolCallbackResponse | PlainMessage<SetToolCallbackResponse> | undefined, b: SetToolCallbackResponse | PlainMessage<SetToolCallbackResponse> | undefined): boolean;
}
//# sourceMappingURL=sdk_bridge_control_service_pb.d.ts.map