import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Duration, Message, proto3 } from "@bufbuild/protobuf";
/**
 * Public SDK error codes returned by the bridge.
 *
 * @generated from enum sdk.v1.SdkErrorCode
 */
export declare enum SdkErrorCode {
    /**
     * @generated from enum value: SDK_ERROR_CODE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: SDK_ERROR_CODE_UNAUTHORIZED = 1;
     */
    UNAUTHORIZED = 1,
    /**
     * @generated from enum value: SDK_ERROR_CODE_API_KEY_NOT_FOUND = 2;
     */
    API_KEY_NOT_FOUND = 2,
    /**
     * @generated from enum value: SDK_ERROR_CODE_PLAN_REQUIRED = 3;
     */
    PLAN_REQUIRED = 3,
    /**
     * @generated from enum value: SDK_ERROR_CODE_ROLE_FORBIDDEN = 4;
     */
    ROLE_FORBIDDEN = 4,
    /**
     * @generated from enum value: SDK_ERROR_CODE_FEATURE_UNAVAILABLE = 5;
     */
    FEATURE_UNAVAILABLE = 5,
    /**
     * @generated from enum value: SDK_ERROR_CODE_AGENT_NOT_FOUND = 6;
     */
    AGENT_NOT_FOUND = 6,
    /**
     * @generated from enum value: SDK_ERROR_CODE_RUN_NOT_FOUND = 7;
     */
    RUN_NOT_FOUND = 7,
    /**
     * @generated from enum value: SDK_ERROR_CODE_VALIDATION_ERROR = 8;
     */
    VALIDATION_ERROR = 8,
    /**
     * @generated from enum value: SDK_ERROR_CODE_INVALID_MODEL = 9;
     */
    INVALID_MODEL = 9,
    /**
     * @generated from enum value: SDK_ERROR_CODE_INVALID_BRANCH_NAME = 10;
     */
    INVALID_BRANCH_NAME = 10,
    /**
     * @generated from enum value: SDK_ERROR_CODE_REPOSITORY_REQUIRED = 11;
     */
    REPOSITORY_REQUIRED = 11,
    /**
     * @generated from enum value: SDK_ERROR_CODE_REPOSITORY_ACCESS = 12;
     */
    REPOSITORY_ACCESS = 12,
    /**
     * @generated from enum value: SDK_ERROR_CODE_PR_RESOLUTION_FAILED = 13;
     */
    PR_RESOLUTION_FAILED = 13,
    /**
     * @generated from enum value: SDK_ERROR_CODE_USAGE_LIMIT_EXCEEDED = 14;
     */
    USAGE_LIMIT_EXCEEDED = 14,
    /**
     * @generated from enum value: SDK_ERROR_CODE_AGENT_BUSY = 15;
     */
    AGENT_BUSY = 15,
    /**
     * @generated from enum value: SDK_ERROR_CODE_AGENT_ARCHIVED = 16;
     */
    AGENT_ARCHIVED = 16,
    /**
     * @generated from enum value: SDK_ERROR_CODE_RUN_NOT_CANCELLABLE = 17;
     */
    RUN_NOT_CANCELLABLE = 17,
    /**
     * @generated from enum value: SDK_ERROR_CODE_RATE_LIMIT_EXCEEDED = 18;
     */
    RATE_LIMIT_EXCEEDED = 18,
    /**
     * @generated from enum value: SDK_ERROR_CODE_UPSTREAM_ERROR = 19;
     */
    UPSTREAM_ERROR = 19,
    /**
     * @generated from enum value: SDK_ERROR_CODE_INTERNAL_ERROR = 20;
     */
    INTERNAL_ERROR = 20,
    /**
     * @generated from enum value: SDK_ERROR_CODE_CLIENT_CANCELLED = 21;
     */
    CLIENT_CANCELLED = 21
}
/**
 * Rate-limit metadata surfaced from Cursor Cloud responses.
 *
 * @generated from message sdk.v1.RateLimitInfo
 */
export declare class RateLimitInfo extends Message<RateLimitInfo> {
    /**
     * @generated from field: optional uint64 limit = 1;
     */
    limit?: bigint;
    /**
     * @generated from field: optional uint64 remaining = 2;
     */
    remaining?: bigint;
    /**
     * @generated from field: optional uint64 reset_epoch_seconds = 3;
     */
    resetEpochSeconds?: bigint;
    constructor(data?: PartialMessage<RateLimitInfo>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RateLimitInfo";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RateLimitInfo;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RateLimitInfo;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RateLimitInfo;
    static equals(a: RateLimitInfo | PlainMessage<RateLimitInfo> | undefined, b: RateLimitInfo | PlainMessage<RateLimitInfo> | undefined): boolean;
}
/**
 * Stable, public error details attached to Connect errors.
 *
 * @generated from message sdk.v1.SdkErrorDetails
 */
export declare class SdkErrorDetails extends Message<SdkErrorDetails> {
    /**
     * Full request ID from Cursor Cloud, if one is available. Never truncate this.
     *
     * @generated from field: optional string request_id = 1;
     */
    requestId?: string;
    /**
     * @generated from field: sdk.v1.SdkErrorCode sdk_error_code = 2;
     */
    sdkErrorCode: SdkErrorCode;
    /**
     * @generated from field: string message = 3;
     */
    message: string;
    /**
     * @generated from field: optional string help_url = 4;
     */
    helpUrl?: string;
    /**
     * @generated from field: optional string provider = 5;
     */
    provider?: string;
    /**
     * @generated from field: optional google.protobuf.Duration retry_after = 6;
     */
    retryAfter?: Duration;
    /**
     * @generated from field: optional sdk.v1.RateLimitInfo rate_limit = 7;
     */
    rateLimit?: RateLimitInfo;
    constructor(data?: PartialMessage<SdkErrorDetails>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkErrorDetails";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkErrorDetails;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkErrorDetails;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkErrorDetails;
    static equals(a: SdkErrorDetails | PlainMessage<SdkErrorDetails> | undefined, b: SdkErrorDetails | PlainMessage<SdkErrorDetails> | undefined): boolean;
}
//# sourceMappingURL=sdk_errors_pb.d.ts.map