import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Struct } from "@bufbuild/protobuf";
/**
 * @generated from message sdk.v1.CallCustomToolRequest
 */
export declare class CallCustomToolRequest extends Message<CallCustomToolRequest> {
    /**
     * @generated from field: string tool_name = 1;
     */
    toolName: string;
    /**
     * @generated from field: google.protobuf.Struct args = 2;
     */
    args?: Struct;
    /**
     * @generated from field: optional string tool_call_id = 3;
     */
    toolCallId?: string;
    /**
     * Agent that owns the tool definition (CreateAgent / ResumeAgent local.custom_tools).
     *
     * @generated from field: string agent_id = 4;
     */
    agentId: string;
    constructor(data?: PartialMessage<CallCustomToolRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CallCustomToolRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CallCustomToolRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CallCustomToolRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CallCustomToolRequest;
    static equals(a: CallCustomToolRequest | PlainMessage<CallCustomToolRequest> | undefined, b: CallCustomToolRequest | PlainMessage<CallCustomToolRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CallCustomToolResponse
 */
export declare class CallCustomToolResponse extends Message<CallCustomToolResponse> {
    /**
     * SDKCustomToolResult wire shape (string, object, or content envelope).
     *
     * @generated from field: google.protobuf.Struct result = 1;
     */
    result?: Struct;
    constructor(data?: PartialMessage<CallCustomToolResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CallCustomToolResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CallCustomToolResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CallCustomToolResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CallCustomToolResponse;
    static equals(a: CallCustomToolResponse | PlainMessage<CallCustomToolResponse> | undefined, b: CallCustomToolResponse | PlainMessage<CallCustomToolResponse> | undefined): boolean;
}
//# sourceMappingURL=sdk_custom_tool_callback_service_pb.d.ts.map