import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Struct } from "@bufbuild/protobuf";
/**
 * Generic LocalAgentStore operation forwarded to a host-owned substore.
 *
 * @generated from message sdk.v1.CallStoreRequest
 */
export declare class CallStoreRequest extends Message<CallStoreRequest> {
    /**
     * Substore being addressed: "agents", "runs", "runEvents", or "checkpoints".
     *
     * @generated from field: string substore = 1;
     */
    substore: string;
    /**
     * Operation on the substore: "get", "create", "update", "delete", "list", or
     * "append" (runEvents only).
     *
     * @generated from field: string method = 2;
     */
    method: string;
    /**
     * Operation input object. Checkpoint blob bytes are base64 strings.
     *
     * @generated from field: google.protobuf.Struct input = 3;
     */
    input?: Struct;
    constructor(data?: PartialMessage<CallStoreRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CallStoreRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CallStoreRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CallStoreRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CallStoreRequest;
    static equals(a: CallStoreRequest | PlainMessage<CallStoreRequest> | undefined, b: CallStoreRequest | PlainMessage<CallStoreRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CallStoreResponse
 */
export declare class CallStoreResponse extends Message<CallStoreResponse> {
    /**
     * Operation output object, or unset for a null result (e.g. a `get` miss or
     * a `delete`). Checkpoint blob bytes are base64 strings.
     *
     * @generated from field: google.protobuf.Struct output = 1;
     */
    output?: Struct;
    constructor(data?: PartialMessage<CallStoreResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CallStoreResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CallStoreResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CallStoreResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CallStoreResponse;
    static equals(a: CallStoreResponse | PlainMessage<CallStoreResponse> | undefined, b: CallStoreResponse | PlainMessage<CallStoreResponse> | undefined): boolean;
}
//# sourceMappingURL=sdk_store_callback_service_pb.d.ts.map