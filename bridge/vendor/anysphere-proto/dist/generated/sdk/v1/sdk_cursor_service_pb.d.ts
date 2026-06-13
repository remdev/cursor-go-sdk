import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { SdkModel, SdkRepository, SdkUser } from "./sdk_messages_pb.js";
/**
 * @generated from message sdk.v1.CursorRequestOptions
 */
export declare class CursorRequestOptions extends Message<CursorRequestOptions> {
    /**
     * @generated from field: string api_key = 1;
     */
    apiKey: string;
    constructor(data?: PartialMessage<CursorRequestOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CursorRequestOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CursorRequestOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CursorRequestOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CursorRequestOptions;
    static equals(a: CursorRequestOptions | PlainMessage<CursorRequestOptions> | undefined, b: CursorRequestOptions | PlainMessage<CursorRequestOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.MeRequest
 */
export declare class MeRequest extends Message<MeRequest> {
    /**
     * @generated from field: sdk.v1.CursorRequestOptions options = 1;
     */
    options?: CursorRequestOptions;
    constructor(data?: PartialMessage<MeRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.MeRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MeRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MeRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MeRequest;
    static equals(a: MeRequest | PlainMessage<MeRequest> | undefined, b: MeRequest | PlainMessage<MeRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.MeResponse
 */
export declare class MeResponse extends Message<MeResponse> {
    /**
     * @generated from field: sdk.v1.SdkUser user = 1;
     */
    user?: SdkUser;
    constructor(data?: PartialMessage<MeResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.MeResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MeResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MeResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MeResponse;
    static equals(a: MeResponse | PlainMessage<MeResponse> | undefined, b: MeResponse | PlainMessage<MeResponse> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ListModelsRequest
 */
export declare class ListModelsRequest extends Message<ListModelsRequest> {
    /**
     * @generated from field: sdk.v1.CursorRequestOptions options = 1;
     */
    options?: CursorRequestOptions;
    constructor(data?: PartialMessage<ListModelsRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ListModelsRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListModelsRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListModelsRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListModelsRequest;
    static equals(a: ListModelsRequest | PlainMessage<ListModelsRequest> | undefined, b: ListModelsRequest | PlainMessage<ListModelsRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ListModelsResponse
 */
export declare class ListModelsResponse extends Message<ListModelsResponse> {
    /**
     * @generated from field: repeated sdk.v1.SdkModel items = 1;
     */
    items: SdkModel[];
    constructor(data?: PartialMessage<ListModelsResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ListModelsResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListModelsResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListModelsResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListModelsResponse;
    static equals(a: ListModelsResponse | PlainMessage<ListModelsResponse> | undefined, b: ListModelsResponse | PlainMessage<ListModelsResponse> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ListRepositoriesRequest
 */
export declare class ListRepositoriesRequest extends Message<ListRepositoriesRequest> {
    /**
     * @generated from field: sdk.v1.CursorRequestOptions options = 1;
     */
    options?: CursorRequestOptions;
    constructor(data?: PartialMessage<ListRepositoriesRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ListRepositoriesRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListRepositoriesRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListRepositoriesRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListRepositoriesRequest;
    static equals(a: ListRepositoriesRequest | PlainMessage<ListRepositoriesRequest> | undefined, b: ListRepositoriesRequest | PlainMessage<ListRepositoriesRequest> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ListRepositoriesResponse
 */
export declare class ListRepositoriesResponse extends Message<ListRepositoriesResponse> {
    /**
     * @generated from field: repeated sdk.v1.SdkRepository items = 1;
     */
    items: SdkRepository[];
    constructor(data?: PartialMessage<ListRepositoriesResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ListRepositoriesResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListRepositoriesResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListRepositoriesResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListRepositoriesResponse;
    static equals(a: ListRepositoriesResponse | PlainMessage<ListRepositoriesResponse> | undefined, b: ListRepositoriesResponse | PlainMessage<ListRepositoriesResponse> | undefined): boolean;
}
//# sourceMappingURL=sdk_cursor_service_pb.d.ts.map