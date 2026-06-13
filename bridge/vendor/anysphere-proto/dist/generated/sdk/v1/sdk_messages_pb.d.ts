import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Struct, Timestamp } from "@bufbuild/protobuf";
/**
 * Runtime selector for SDK APIs that can route to local or cloud agents.
 *
 * @generated from enum sdk.v1.Runtime
 */
export declare enum Runtime {
    /**
     * @generated from enum value: RUNTIME_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: RUNTIME_AUTO = 1;
     */
    AUTO = 1,
    /**
     * @generated from enum value: RUNTIME_LOCAL = 2;
     */
    LOCAL = 2,
    /**
     * @generated from enum value: RUNTIME_CLOUD = 3;
     */
    CLOUD = 3
}
/**
 * Lifecycle status used for run snapshots and stream status messages.
 *
 * @generated from enum sdk.v1.RunLifecycleStatus
 */
export declare enum RunLifecycleStatus {
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_CREATING = 1;
     */
    CREATING = 1,
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_RUNNING = 2;
     */
    RUNNING = 2,
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_FINISHED = 3;
     */
    FINISHED = 3,
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_ERROR = 4;
     */
    ERROR = 4,
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_CANCELLED = 5;
     */
    CANCELLED = 5,
    /**
     * @generated from enum value: RUN_LIFECYCLE_STATUS_EXPIRED = 6;
     */
    EXPIRED = 6
}
/**
 * Agent lifecycle status returned by list/get agent calls.
 *
 * @generated from enum sdk.v1.AgentInfoStatus
 */
export declare enum AgentInfoStatus {
    /**
     * @generated from enum value: AGENT_INFO_STATUS_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: AGENT_INFO_STATUS_RUNNING = 1;
     */
    RUNNING = 1,
    /**
     * @generated from enum value: AGENT_INFO_STATUS_FINISHED = 2;
     */
    FINISHED = 2,
    /**
     * @generated from enum value: AGENT_INFO_STATUS_ERROR = 3;
     */
    ERROR = 3
}
/**
 * Cursor setting sources accepted by local agent options.
 *
 * @generated from enum sdk.v1.SettingSource
 */
export declare enum SettingSource {
    /**
     * @generated from enum value: SETTING_SOURCE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: SETTING_SOURCE_PROJECT = 1;
     */
    PROJECT = 1,
    /**
     * @generated from enum value: SETTING_SOURCE_USER = 2;
     */
    USER = 2,
    /**
     * @generated from enum value: SETTING_SOURCE_TEAM = 3;
     */
    TEAM = 3,
    /**
     * @generated from enum value: SETTING_SOURCE_MDM = 4;
     */
    MDM = 4,
    /**
     * @generated from enum value: SETTING_SOURCE_PLUGINS = 5;
     */
    PLUGINS = 5,
    /**
     * @generated from enum value: SETTING_SOURCE_ALL = 6;
     */
    ALL = 6
}
/**
 * Cloud execution environment type.
 *
 * @generated from enum sdk.v1.CloudEnvironmentType
 */
export declare enum CloudEnvironmentType {
    /**
     * @generated from enum value: CLOUD_ENVIRONMENT_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: CLOUD_ENVIRONMENT_TYPE_CLOUD = 1;
     */
    CLOUD = 1,
    /**
     * @generated from enum value: CLOUD_ENVIRONMENT_TYPE_POOL = 2;
     */
    POOL = 2,
    /**
     * @generated from enum value: CLOUD_ENVIRONMENT_TYPE_MACHINE = 3;
     */
    MACHINE = 3
}
/**
 * MCP transport used for HTTP-style MCP servers.
 *
 * @generated from enum sdk.v1.HttpMcpTransportType
 */
export declare enum HttpMcpTransportType {
    /**
     * @generated from enum value: HTTP_MCP_TRANSPORT_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: HTTP_MCP_TRANSPORT_TYPE_HTTP = 1;
     */
    HTTP = 1,
    /**
     * @generated from enum value: HTTP_MCP_TRANSPORT_TYPE_SSE = 2;
     */
    SSE = 2
}
/**
 * SDK conversation mode for agent creation and per-send overrides.
 *
 * @generated from enum sdk.v1.AgentModeOption
 */
export declare enum AgentModeOption {
    /**
     * @generated from enum value: AGENT_MODE_OPTION_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from enum value: AGENT_MODE_OPTION_AGENT = 1;
     */
    AGENT = 1,
    /**
     * @generated from enum value: AGENT_MODE_OPTION_PLAN = 2;
     */
    PLAN = 2
}
/**
 * @generated from message sdk.v1.ModelSelection
 */
export declare class ModelSelection extends Message<ModelSelection> {
    /**
     * @generated from field: string id = 1;
     */
    id: string;
    /**
     * @generated from field: repeated sdk.v1.ModelParameterValue params = 2;
     */
    params: ModelParameterValue[];
    constructor(data?: PartialMessage<ModelSelection>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ModelSelection";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ModelSelection;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ModelSelection;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ModelSelection;
    static equals(a: ModelSelection | PlainMessage<ModelSelection> | undefined, b: ModelSelection | PlainMessage<ModelSelection> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ModelParameterValue
 */
export declare class ModelParameterValue extends Message<ModelParameterValue> {
    /**
     * @generated from field: string id = 1;
     */
    id: string;
    /**
     * @generated from field: string value = 2;
     */
    value: string;
    constructor(data?: PartialMessage<ModelParameterValue>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ModelParameterValue";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ModelParameterValue;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ModelParameterValue;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ModelParameterValue;
    static equals(a: ModelParameterValue | PlainMessage<ModelParameterValue> | undefined, b: ModelParameterValue | PlainMessage<ModelParameterValue> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ModelParameterDefinition
 */
export declare class ModelParameterDefinition extends Message<ModelParameterDefinition> {
    /**
     * @generated from field: string id = 1;
     */
    id: string;
    /**
     * @generated from field: string display_name = 2;
     */
    displayName: string;
    /**
     * @generated from field: repeated sdk.v1.ModelParameterDefinitionValue values = 3;
     */
    values: ModelParameterDefinitionValue[];
    constructor(data?: PartialMessage<ModelParameterDefinition>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ModelParameterDefinition";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ModelParameterDefinition;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ModelParameterDefinition;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ModelParameterDefinition;
    static equals(a: ModelParameterDefinition | PlainMessage<ModelParameterDefinition> | undefined, b: ModelParameterDefinition | PlainMessage<ModelParameterDefinition> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ModelParameterDefinitionValue
 */
export declare class ModelParameterDefinitionValue extends Message<ModelParameterDefinitionValue> {
    /**
     * @generated from field: string value = 1;
     */
    value: string;
    /**
     * @generated from field: string display_name = 2;
     */
    displayName: string;
    constructor(data?: PartialMessage<ModelParameterDefinitionValue>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ModelParameterDefinitionValue";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ModelParameterDefinitionValue;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ModelParameterDefinitionValue;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ModelParameterDefinitionValue;
    static equals(a: ModelParameterDefinitionValue | PlainMessage<ModelParameterDefinitionValue> | undefined, b: ModelParameterDefinitionValue | PlainMessage<ModelParameterDefinitionValue> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ModelVariant
 */
export declare class ModelVariant extends Message<ModelVariant> {
    /**
     * @generated from field: repeated sdk.v1.ModelParameterValue params = 1;
     */
    params: ModelParameterValue[];
    /**
     * @generated from field: string display_name = 2;
     */
    displayName: string;
    /**
     * @generated from field: string description = 3;
     */
    description: string;
    /**
     * @generated from field: bool is_default = 4;
     */
    isDefault: boolean;
    constructor(data?: PartialMessage<ModelVariant>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ModelVariant";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ModelVariant;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ModelVariant;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ModelVariant;
    static equals(a: ModelVariant | PlainMessage<ModelVariant> | undefined, b: ModelVariant | PlainMessage<ModelVariant> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkModel
 */
export declare class SdkModel extends Message<SdkModel> {
    /**
     * @generated from field: string id = 1;
     */
    id: string;
    /**
     * @generated from field: string display_name = 2;
     */
    displayName: string;
    /**
     * @generated from field: string description = 3;
     */
    description: string;
    /**
     * @generated from field: repeated sdk.v1.ModelParameterDefinition parameters = 4;
     */
    parameters: ModelParameterDefinition[];
    /**
     * @generated from field: repeated sdk.v1.ModelVariant variants = 5;
     */
    variants: ModelVariant[];
    constructor(data?: PartialMessage<SdkModel>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkModel";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkModel;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkModel;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkModel;
    static equals(a: SdkModel | PlainMessage<SdkModel> | undefined, b: SdkModel | PlainMessage<SdkModel> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkUser
 */
export declare class SdkUser extends Message<SdkUser> {
    /**
     * @generated from field: string api_key_name = 1;
     */
    apiKeyName: string;
    /**
     * @generated from field: uint64 user_id = 2;
     */
    userId: bigint;
    /**
     * @generated from field: string user_email = 3;
     */
    userEmail: string;
    /**
     * @generated from field: string user_first_name = 4;
     */
    userFirstName: string;
    /**
     * @generated from field: string user_last_name = 5;
     */
    userLastName: string;
    /**
     * @generated from field: string created_at = 6;
     */
    createdAt: string;
    constructor(data?: PartialMessage<SdkUser>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkUser";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkUser;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkUser;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkUser;
    static equals(a: SdkUser | PlainMessage<SdkUser> | undefined, b: SdkUser | PlainMessage<SdkUser> | undefined): boolean;
}
/**
 * Current @cursor/sdk repository shape. Rich metadata can be added later.
 *
 * @generated from message sdk.v1.SdkRepository
 */
export declare class SdkRepository extends Message<SdkRepository> {
    /**
     * @generated from field: string url = 1;
     */
    url: string;
    constructor(data?: PartialMessage<SdkRepository>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkRepository";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkRepository;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkRepository;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkRepository;
    static equals(a: SdkRepository | PlainMessage<SdkRepository> | undefined, b: SdkRepository | PlainMessage<SdkRepository> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.UserMessage
 */
export declare class UserMessage extends Message<UserMessage> {
    /**
     * @generated from field: string text = 1;
     */
    text: string;
    /**
     * @generated from field: repeated sdk.v1.SdkImage images = 2;
     */
    images: SdkImage[];
    constructor(data?: PartialMessage<UserMessage>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.UserMessage";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UserMessage;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UserMessage;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UserMessage;
    static equals(a: UserMessage | PlainMessage<UserMessage> | undefined, b: UserMessage | PlainMessage<UserMessage> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkImage
 */
export declare class SdkImage extends Message<SdkImage> {
    /**
     * @generated from oneof sdk.v1.SdkImage.source
     */
    source: {
        /**
         * @generated from field: sdk.v1.SdkImageUrl url = 1;
         */
        value: SdkImageUrl;
        case: "url";
    } | {
        /**
         * @generated from field: sdk.v1.SdkImageData data = 2;
         */
        value: SdkImageData;
        case: "data";
    } | {
        case: undefined;
        value?: undefined;
    };
    /**
     * @generated from field: sdk.v1.SdkImageDimension dimension = 3;
     */
    dimension?: SdkImageDimension;
    constructor(data?: PartialMessage<SdkImage>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkImage";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkImage;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkImage;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkImage;
    static equals(a: SdkImage | PlainMessage<SdkImage> | undefined, b: SdkImage | PlainMessage<SdkImage> | undefined): boolean;
}
/**
 * URL images are accepted only by cloud-routed agents in v1.
 *
 * @generated from message sdk.v1.SdkImageUrl
 */
export declare class SdkImageUrl extends Message<SdkImageUrl> {
    /**
     * @generated from field: string url = 1;
     */
    url: string;
    constructor(data?: PartialMessage<SdkImageUrl>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkImageUrl";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkImageUrl;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkImageUrl;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkImageUrl;
    static equals(a: SdkImageUrl | PlainMessage<SdkImageUrl> | undefined, b: SdkImageUrl | PlainMessage<SdkImageUrl> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkImageData
 */
export declare class SdkImageData extends Message<SdkImageData> {
    /**
     * @generated from field: string data = 1;
     */
    data: string;
    /**
     * @generated from field: string mime_type = 2;
     */
    mimeType: string;
    constructor(data?: PartialMessage<SdkImageData>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkImageData";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkImageData;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkImageData;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkImageData;
    static equals(a: SdkImageData | PlainMessage<SdkImageData> | undefined, b: SdkImageData | PlainMessage<SdkImageData> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkImageDimension
 */
export declare class SdkImageDimension extends Message<SdkImageDimension> {
    /**
     * @generated from field: uint32 width = 1;
     */
    width: number;
    /**
     * @generated from field: uint32 height = 2;
     */
    height: number;
    constructor(data?: PartialMessage<SdkImageDimension>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkImageDimension";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkImageDimension;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkImageDimension;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkImageDimension;
    static equals(a: SdkImageDimension | PlainMessage<SdkImageDimension> | undefined, b: SdkImageDimension | PlainMessage<SdkImageDimension> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.AgentOptions
 */
export declare class AgentOptions extends Message<AgentOptions> {
    /**
     * @generated from field: sdk.v1.ModelSelection model = 1;
     */
    model?: ModelSelection;
    /**
     * @generated from field: string api_key = 2;
     */
    apiKey: string;
    /**
     * @generated from field: string name = 3;
     */
    name: string;
    /**
     * @generated from field: sdk.v1.LocalAgentOptions local = 4;
     */
    local?: LocalAgentOptions;
    /**
     * @generated from field: sdk.v1.CloudAgentOptions cloud = 5;
     */
    cloud?: CloudAgentOptions;
    /**
     * @generated from field: map<string, sdk.v1.McpServerConfig> mcp_servers = 6;
     */
    mcpServers: {
        [key: string]: McpServerConfig;
    };
    /**
     * @generated from field: map<string, sdk.v1.AgentDefinition> agents = 7;
     */
    agents: {
        [key: string]: AgentDefinition;
    };
    /**
     * @generated from field: string agent_id = 8;
     */
    agentId: string;
    /**
     * Initial conversation mode for the agent.
     *
     * @generated from field: sdk.v1.AgentModeOption mode = 9;
     */
    mode: AgentModeOption;
    constructor(data?: PartialMessage<AgentOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.AgentOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AgentOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AgentOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AgentOptions;
    static equals(a: AgentOptions | PlainMessage<AgentOptions> | undefined, b: AgentOptions | PlainMessage<AgentOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SendOptions
 */
export declare class SendOptions extends Message<SendOptions> {
    /**
     * @generated from field: sdk.v1.ModelSelection model = 1;
     */
    model?: ModelSelection;
    /**
     * @generated from field: map<string, sdk.v1.McpServerConfig> mcp_servers = 2;
     */
    mcpServers: {
        [key: string]: McpServerConfig;
    };
    /**
     * @generated from field: sdk.v1.LocalSendOptions local = 3;
     */
    local?: LocalSendOptions;
    /**
     * When true, the bridge forwards raw InteractionUpdate callbacks from
     * @cursor/sdk as RunStreamMessage.interaction_update events.
     *
     * @generated from field: bool enable_deltas = 4;
     */
    enableDeltas: boolean;
    /**
     * When true, the bridge forwards completed ConversationStep callbacks from
     * @cursor/sdk as RunStreamMessage.step events.
     *
     * @generated from field: bool enable_steps = 5;
     */
    enableSteps: boolean;
    /**
     * Per-send conversation mode override.
     *
     * @generated from field: sdk.v1.AgentModeOption mode = 6;
     */
    mode: AgentModeOption;
    constructor(data?: PartialMessage<SendOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SendOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SendOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SendOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SendOptions;
    static equals(a: SendOptions | PlainMessage<SendOptions> | undefined, b: SendOptions | PlainMessage<SendOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.LocalSendOptions
 */
export declare class LocalSendOptions extends Message<LocalSendOptions> {
    /**
     * Local-only force send; cloud-routed agents reject this field.
     *
     * @generated from field: optional bool force = 1;
     */
    force?: boolean;
    constructor(data?: PartialMessage<LocalSendOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.LocalSendOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LocalSendOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LocalSendOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LocalSendOptions;
    static equals(a: LocalSendOptions | PlainMessage<LocalSendOptions> | undefined, b: LocalSendOptions | PlainMessage<LocalSendOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.LocalAgentStoreConfig
 */
export declare class LocalAgentStoreConfig extends Message<LocalAgentStoreConfig> {
    /**
     * Store topology: "sqlite" (bridge-owned SQLite; also the default when unset),
     * "jsonl" (bridge-owned JSONL files), or "custom" (fully host-owned via the
     * SdkStoreCallbackService).
     *
     * @generated from field: string type = 1;
     */
    type: string;
    /**
     * Required for "jsonl".
     *
     * @generated from field: string root_dir = 2;
     */
    rootDir: string;
    constructor(data?: PartialMessage<LocalAgentStoreConfig>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.LocalAgentStoreConfig";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LocalAgentStoreConfig;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LocalAgentStoreConfig;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LocalAgentStoreConfig;
    static equals(a: LocalAgentStoreConfig | PlainMessage<LocalAgentStoreConfig> | undefined, b: LocalAgentStoreConfig | PlainMessage<LocalAgentStoreConfig> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CustomToolDefinition
 */
export declare class CustomToolDefinition extends Message<CustomToolDefinition> {
    /**
     * @generated from field: optional string description = 1;
     */
    description?: string;
    /**
     * @generated from field: google.protobuf.Struct input_schema = 2;
     */
    inputSchema?: Struct;
    constructor(data?: PartialMessage<CustomToolDefinition>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CustomToolDefinition";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CustomToolDefinition;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CustomToolDefinition;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CustomToolDefinition;
    static equals(a: CustomToolDefinition | PlainMessage<CustomToolDefinition> | undefined, b: CustomToolDefinition | PlainMessage<CustomToolDefinition> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.LocalAgentOptions
 */
export declare class LocalAgentOptions extends Message<LocalAgentOptions> {
    /**
     * @generated from field: repeated string cwd = 1;
     */
    cwd: string[];
    /**
     * @generated from field: repeated sdk.v1.SettingSource setting_sources = 2;
     */
    settingSources: SettingSource[];
    /**
     * @generated from field: sdk.v1.SandboxOptions sandbox_options = 3;
     */
    sandboxOptions?: SandboxOptions;
    /**
     * @generated from field: sdk.v1.LocalAgentStoreConfig store = 4;
     */
    store?: LocalAgentStoreConfig;
    /**
     * When true, local runs use Smart Auto Review (classifier-backed Auto) for tool calls.
     *
     * @generated from field: optional bool auto_review = 5;
     */
    autoReview?: boolean;
    /**
     * Agent-level custom tools (metadata only). Execution is host-owned via
     * SdkCustomToolCallbackService when using the Python SDK bridge.
     *
     * @generated from field: map<string, sdk.v1.CustomToolDefinition> custom_tools = 6;
     */
    customTools: {
        [key: string]: CustomToolDefinition;
    };
    constructor(data?: PartialMessage<LocalAgentOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.LocalAgentOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LocalAgentOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LocalAgentOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LocalAgentOptions;
    static equals(a: LocalAgentOptions | PlainMessage<LocalAgentOptions> | undefined, b: LocalAgentOptions | PlainMessage<LocalAgentOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SandboxOptions
 */
export declare class SandboxOptions extends Message<SandboxOptions> {
    /**
     * @generated from field: optional bool enabled = 1;
     */
    enabled?: boolean;
    constructor(data?: PartialMessage<SandboxOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SandboxOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SandboxOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SandboxOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SandboxOptions;
    static equals(a: SandboxOptions | PlainMessage<SandboxOptions> | undefined, b: SandboxOptions | PlainMessage<SandboxOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CloudAgentOptions
 */
export declare class CloudAgentOptions extends Message<CloudAgentOptions> {
    /**
     * @generated from field: sdk.v1.CloudEnvironment env = 1;
     */
    env?: CloudEnvironment;
    /**
     * @generated from field: repeated sdk.v1.CloudRepository repos = 2;
     */
    repos: CloudRepository[];
    /**
     * @generated from field: optional bool work_on_current_branch = 3;
     */
    workOnCurrentBranch?: boolean;
    /**
     * @generated from field: optional bool auto_create_pr = 4;
     */
    autoCreatePr?: boolean;
    /**
     * @generated from field: optional bool skip_reviewer_request = 5;
     */
    skipReviewerRequest?: boolean;
    /**
     * @generated from field: map<string, string> env_vars = 6;
     */
    envVars: {
        [key: string]: string;
    };
    constructor(data?: PartialMessage<CloudAgentOptions>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CloudAgentOptions";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CloudAgentOptions;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CloudAgentOptions;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CloudAgentOptions;
    static equals(a: CloudAgentOptions | PlainMessage<CloudAgentOptions> | undefined, b: CloudAgentOptions | PlainMessage<CloudAgentOptions> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CloudEnvironment
 */
export declare class CloudEnvironment extends Message<CloudEnvironment> {
    /**
     * @generated from field: sdk.v1.CloudEnvironmentType type = 1;
     */
    type: CloudEnvironmentType;
    /**
     * @generated from field: string name = 2;
     */
    name: string;
    constructor(data?: PartialMessage<CloudEnvironment>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CloudEnvironment";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CloudEnvironment;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CloudEnvironment;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CloudEnvironment;
    static equals(a: CloudEnvironment | PlainMessage<CloudEnvironment> | undefined, b: CloudEnvironment | PlainMessage<CloudEnvironment> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CloudRepository
 */
export declare class CloudRepository extends Message<CloudRepository> {
    /**
     * @generated from field: string url = 1;
     */
    url: string;
    /**
     * @generated from field: string starting_ref = 2;
     */
    startingRef: string;
    /**
     * @generated from field: string pr_url = 3;
     */
    prUrl: string;
    constructor(data?: PartialMessage<CloudRepository>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CloudRepository";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CloudRepository;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CloudRepository;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CloudRepository;
    static equals(a: CloudRepository | PlainMessage<CloudRepository> | undefined, b: CloudRepository | PlainMessage<CloudRepository> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.McpServerConfig
 */
export declare class McpServerConfig extends Message<McpServerConfig> {
    /**
     * @generated from oneof sdk.v1.McpServerConfig.config
     */
    config: {
        /**
         * @generated from field: sdk.v1.StdioMcpServerConfig stdio = 1;
         */
        value: StdioMcpServerConfig;
        case: "stdio";
    } | {
        /**
         * @generated from field: sdk.v1.HttpMcpServerConfig http = 2;
         */
        value: HttpMcpServerConfig;
        case: "http";
    } | {
        case: undefined;
        value?: undefined;
    };
    constructor(data?: PartialMessage<McpServerConfig>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.McpServerConfig";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): McpServerConfig;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): McpServerConfig;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): McpServerConfig;
    static equals(a: McpServerConfig | PlainMessage<McpServerConfig> | undefined, b: McpServerConfig | PlainMessage<McpServerConfig> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.StdioMcpServerConfig
 */
export declare class StdioMcpServerConfig extends Message<StdioMcpServerConfig> {
    /**
     * @generated from field: string command = 1;
     */
    command: string;
    /**
     * @generated from field: repeated string args = 2;
     */
    args: string[];
    /**
     * @generated from field: map<string, string> env = 3;
     */
    env: {
        [key: string]: string;
    };
    /**
     * @generated from field: string cwd = 4;
     */
    cwd: string;
    constructor(data?: PartialMessage<StdioMcpServerConfig>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.StdioMcpServerConfig";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StdioMcpServerConfig;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StdioMcpServerConfig;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StdioMcpServerConfig;
    static equals(a: StdioMcpServerConfig | PlainMessage<StdioMcpServerConfig> | undefined, b: StdioMcpServerConfig | PlainMessage<StdioMcpServerConfig> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.HttpMcpServerConfig
 */
export declare class HttpMcpServerConfig extends Message<HttpMcpServerConfig> {
    /**
     * @generated from field: sdk.v1.HttpMcpTransportType type = 1;
     */
    type: HttpMcpTransportType;
    /**
     * @generated from field: string url = 2;
     */
    url: string;
    /**
     * @generated from field: map<string, string> headers = 3;
     */
    headers: {
        [key: string]: string;
    };
    /**
     * @generated from field: sdk.v1.McpAuthConfig auth = 4;
     */
    auth?: McpAuthConfig;
    constructor(data?: PartialMessage<HttpMcpServerConfig>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.HttpMcpServerConfig";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): HttpMcpServerConfig;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): HttpMcpServerConfig;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): HttpMcpServerConfig;
    static equals(a: HttpMcpServerConfig | PlainMessage<HttpMcpServerConfig> | undefined, b: HttpMcpServerConfig | PlainMessage<HttpMcpServerConfig> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.McpAuthConfig
 */
export declare class McpAuthConfig extends Message<McpAuthConfig> {
    /**
     * @generated from field: string client_id = 1;
     */
    clientId: string;
    /**
     * @generated from field: string client_secret = 2;
     */
    clientSecret: string;
    /**
     * @generated from field: repeated string scopes = 3;
     */
    scopes: string[];
    constructor(data?: PartialMessage<McpAuthConfig>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.McpAuthConfig";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): McpAuthConfig;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): McpAuthConfig;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): McpAuthConfig;
    static equals(a: McpAuthConfig | PlainMessage<McpAuthConfig> | undefined, b: McpAuthConfig | PlainMessage<McpAuthConfig> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.AgentDefinition
 */
export declare class AgentDefinition extends Message<AgentDefinition> {
    /**
     * @generated from field: string description = 1;
     */
    description: string;
    /**
     * @generated from field: string prompt = 2;
     */
    prompt: string;
    /**
     * @generated from field: sdk.v1.ModelSelection model = 3;
     */
    model?: ModelSelection;
    /**
     * @generated from field: optional bool inherit_model = 4;
     */
    inheritModel?: boolean;
    /**
     * @generated from field: repeated sdk.v1.AgentDefinitionMcpServer mcp_servers = 5;
     */
    mcpServers: AgentDefinitionMcpServer[];
    constructor(data?: PartialMessage<AgentDefinition>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.AgentDefinition";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AgentDefinition;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AgentDefinition;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AgentDefinition;
    static equals(a: AgentDefinition | PlainMessage<AgentDefinition> | undefined, b: AgentDefinition | PlainMessage<AgentDefinition> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.AgentDefinitionMcpServer
 */
export declare class AgentDefinitionMcpServer extends Message<AgentDefinitionMcpServer> {
    /**
     * @generated from oneof sdk.v1.AgentDefinitionMcpServer.value
     */
    value: {
        /**
         * @generated from field: string name = 1;
         */
        value: string;
        case: "name";
    } | {
        /**
         * @generated from field: sdk.v1.McpServerConfig inline_config = 2;
         */
        value: McpServerConfig;
        case: "inlineConfig";
    } | {
        case: undefined;
        value?: undefined;
    };
    constructor(data?: PartialMessage<AgentDefinitionMcpServer>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.AgentDefinitionMcpServer";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AgentDefinitionMcpServer;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AgentDefinitionMcpServer;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AgentDefinitionMcpServer;
    static equals(a: AgentDefinitionMcpServer | PlainMessage<AgentDefinitionMcpServer> | undefined, b: AgentDefinitionMcpServer | PlainMessage<AgentDefinitionMcpServer> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkAgentInfo
 */
export declare class SdkAgentInfo extends Message<SdkAgentInfo> {
    /**
     * @generated from field: string agent_id = 1;
     */
    agentId: string;
    /**
     * @generated from field: string name = 2;
     */
    name: string;
    /**
     * @generated from field: string summary = 3;
     */
    summary: string;
    /**
     * @generated from field: google.protobuf.Timestamp last_modified = 4;
     */
    lastModified?: Timestamp;
    /**
     * @generated from field: sdk.v1.AgentInfoStatus status = 5;
     */
    status: AgentInfoStatus;
    /**
     * @generated from field: google.protobuf.Timestamp created_at = 6;
     */
    createdAt?: Timestamp;
    /**
     * @generated from field: bool archived = 7;
     */
    archived: boolean;
    /**
     * @generated from oneof sdk.v1.SdkAgentInfo.runtime_info
     */
    runtimeInfo: {
        /**
         * @generated from field: sdk.v1.LocalAgentInfo local = 8;
         */
        value: LocalAgentInfo;
        case: "local";
    } | {
        /**
         * @generated from field: sdk.v1.CloudAgentInfo cloud = 9;
         */
        value: CloudAgentInfo;
        case: "cloud";
    } | {
        case: undefined;
        value?: undefined;
    };
    constructor(data?: PartialMessage<SdkAgentInfo>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkAgentInfo";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkAgentInfo;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkAgentInfo;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkAgentInfo;
    static equals(a: SdkAgentInfo | PlainMessage<SdkAgentInfo> | undefined, b: SdkAgentInfo | PlainMessage<SdkAgentInfo> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.LocalAgentInfo
 */
export declare class LocalAgentInfo extends Message<LocalAgentInfo> {
    /**
     * @generated from field: string cwd = 1;
     */
    cwd: string;
    constructor(data?: PartialMessage<LocalAgentInfo>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.LocalAgentInfo";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LocalAgentInfo;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LocalAgentInfo;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LocalAgentInfo;
    static equals(a: LocalAgentInfo | PlainMessage<LocalAgentInfo> | undefined, b: LocalAgentInfo | PlainMessage<LocalAgentInfo> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.CloudAgentInfo
 */
export declare class CloudAgentInfo extends Message<CloudAgentInfo> {
    /**
     * @generated from field: sdk.v1.CloudEnvironment env = 1;
     */
    env?: CloudEnvironment;
    /**
     * @generated from field: repeated string repos = 2;
     */
    repos: string[];
    constructor(data?: PartialMessage<CloudAgentInfo>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.CloudAgentInfo";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CloudAgentInfo;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CloudAgentInfo;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CloudAgentInfo;
    static equals(a: CloudAgentInfo | PlainMessage<CloudAgentInfo> | undefined, b: CloudAgentInfo | PlainMessage<CloudAgentInfo> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.RunSnapshot
 */
export declare class RunSnapshot extends Message<RunSnapshot> {
    /**
     * @generated from field: string run_id = 1;
     */
    runId: string;
    /**
     * @generated from field: string agent_id = 2;
     */
    agentId: string;
    /**
     * @generated from field: sdk.v1.RunLifecycleStatus status = 3;
     */
    status: RunLifecycleStatus;
    /**
     * @generated from field: string result = 4;
     */
    result: string;
    /**
     * @generated from field: sdk.v1.ModelSelection model = 5;
     */
    model?: ModelSelection;
    /**
     * @generated from field: uint64 duration_ms = 6;
     */
    durationMs: bigint;
    /**
     * @generated from field: sdk.v1.RunGitInfo git = 7;
     */
    git?: RunGitInfo;
    /**
     * @generated from field: google.protobuf.Timestamp created_at = 8;
     */
    createdAt?: Timestamp;
    constructor(data?: PartialMessage<RunSnapshot>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunSnapshot";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunSnapshot;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunSnapshot;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunSnapshot;
    static equals(a: RunSnapshot | PlainMessage<RunSnapshot> | undefined, b: RunSnapshot | PlainMessage<RunSnapshot> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.RunResult
 */
export declare class RunResult extends Message<RunResult> {
    /**
     * @generated from field: string run_id = 1;
     */
    runId: string;
    /**
     * @generated from field: string agent_id = 2;
     */
    agentId: string;
    /**
     * @generated from field: sdk.v1.RunLifecycleStatus status = 3;
     */
    status: RunLifecycleStatus;
    /**
     * @generated from field: string result = 4;
     */
    result: string;
    /**
     * @generated from field: sdk.v1.ModelSelection model = 5;
     */
    model?: ModelSelection;
    /**
     * @generated from field: uint64 duration_ms = 6;
     */
    durationMs: bigint;
    /**
     * @generated from field: sdk.v1.RunGitInfo git = 7;
     */
    git?: RunGitInfo;
    /**
     * @generated from field: google.protobuf.Timestamp created_at = 8;
     */
    createdAt?: Timestamp;
    constructor(data?: PartialMessage<RunResult>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunResult";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunResult;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunResult;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunResult;
    static equals(a: RunResult | PlainMessage<RunResult> | undefined, b: RunResult | PlainMessage<RunResult> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.RunGitInfo
 */
export declare class RunGitInfo extends Message<RunGitInfo> {
    /**
     * @generated from field: repeated sdk.v1.RunGitBranchInfo branches = 1;
     */
    branches: RunGitBranchInfo[];
    constructor(data?: PartialMessage<RunGitInfo>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunGitInfo";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunGitInfo;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunGitInfo;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunGitInfo;
    static equals(a: RunGitInfo | PlainMessage<RunGitInfo> | undefined, b: RunGitInfo | PlainMessage<RunGitInfo> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.RunGitBranchInfo
 */
export declare class RunGitBranchInfo extends Message<RunGitBranchInfo> {
    /**
     * @generated from field: string repo_url = 1;
     */
    repoUrl: string;
    /**
     * @generated from field: string branch = 2;
     */
    branch: string;
    /**
     * @generated from field: string pr_url = 3;
     */
    prUrl: string;
    constructor(data?: PartialMessage<RunGitBranchInfo>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunGitBranchInfo";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunGitBranchInfo;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunGitBranchInfo;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunGitBranchInfo;
    static equals(a: RunGitBranchInfo | PlainMessage<RunGitBranchInfo> | undefined, b: RunGitBranchInfo | PlainMessage<RunGitBranchInfo> | undefined): boolean;
}
/**
 * Stream message envelope for live Send and durable ObserveRun streams.
 *
 * @generated from message sdk.v1.RunStreamMessage
 */
export declare class RunStreamMessage extends Message<RunStreamMessage> {
    /**
     * @generated from oneof sdk.v1.RunStreamMessage.envelope
     */
    envelope: {
        /**
         * @generated from field: sdk.v1.SdkMessage sdk_message = 1;
         */
        value: SdkMessage;
        case: "sdkMessage";
    } | {
        /**
         * @generated from field: sdk.v1.RunStreamResult result = 2;
         */
        value: RunStreamResult;
        case: "result";
    } | {
        /**
         * @generated from field: sdk.v1.RunStreamDone done = 3;
         */
        value: RunStreamDone;
        case: "done";
    } | {
        /**
         * @generated from field: sdk.v1.InteractionUpdate interaction_update = 5;
         */
        value: InteractionUpdate;
        case: "interactionUpdate";
    } | {
        /**
         * @generated from field: sdk.v1.ConversationStep step = 6;
         */
        value: ConversationStep;
        case: "step";
    } | {
        case: undefined;
        value?: undefined;
    };
    /**
     * Opaque exclusive resume token, scoped to a single run.
     *
     * @generated from field: optional string offset = 4;
     */
    offset?: string;
    constructor(data?: PartialMessage<RunStreamMessage>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunStreamMessage";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunStreamMessage;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunStreamMessage;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunStreamMessage;
    static equals(a: RunStreamMessage | PlainMessage<RunStreamMessage> | undefined, b: RunStreamMessage | PlainMessage<RunStreamMessage> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.RunStreamResult
 */
export declare class RunStreamResult extends Message<RunStreamResult> {
    /**
     * @generated from field: string agent_id = 1;
     */
    agentId: string;
    /**
     * @generated from field: string run_id = 2;
     */
    runId: string;
    /**
     * @generated from field: sdk.v1.RunLifecycleStatus status = 3;
     */
    status: RunLifecycleStatus;
    /**
     * @generated from field: optional string error_code = 4;
     */
    errorCode?: string;
    /**
     * @generated from field: sdk.v1.RunResult result = 5;
     */
    result?: RunResult;
    constructor(data?: PartialMessage<RunStreamResult>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunStreamResult";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunStreamResult;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunStreamResult;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunStreamResult;
    static equals(a: RunStreamResult | PlainMessage<RunStreamResult> | undefined, b: RunStreamResult | PlainMessage<RunStreamResult> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.RunStreamDone
 */
export declare class RunStreamDone extends Message<RunStreamDone> {
    /**
     * @generated from field: string agent_id = 1;
     */
    agentId: string;
    /**
     * @generated from field: string run_id = 2;
     */
    runId: string;
    constructor(data?: PartialMessage<RunStreamDone>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.RunStreamDone";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunStreamDone;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunStreamDone;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunStreamDone;
    static equals(a: RunStreamDone | PlainMessage<RunStreamDone> | undefined, b: RunStreamDone | PlainMessage<RunStreamDone> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.InteractionUpdate
 */
export declare class InteractionUpdate extends Message<InteractionUpdate> {
    /**
     * @generated from field: string type = 1;
     */
    type: string;
    /**
     * @generated from field: google.protobuf.Struct update = 2;
     */
    update?: Struct;
    constructor(data?: PartialMessage<InteractionUpdate>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.InteractionUpdate";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): InteractionUpdate;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): InteractionUpdate;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): InteractionUpdate;
    static equals(a: InteractionUpdate | PlainMessage<InteractionUpdate> | undefined, b: InteractionUpdate | PlainMessage<InteractionUpdate> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.ConversationStep
 */
export declare class ConversationStep extends Message<ConversationStep> {
    /**
     * @generated from field: string type = 1;
     */
    type: string;
    /**
     * @generated from field: google.protobuf.Struct step = 2;
     */
    step?: Struct;
    constructor(data?: PartialMessage<ConversationStep>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.ConversationStep";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ConversationStep;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ConversationStep;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ConversationStep;
    static equals(a: ConversationStep | PlainMessage<ConversationStep> | undefined, b: ConversationStep | PlainMessage<ConversationStep> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkMessage
 */
export declare class SdkMessage extends Message<SdkMessage> {
    /**
     * @generated from field: string type = 1;
     */
    type: string;
    /**
     * @generated from field: google.protobuf.Struct message = 2;
     */
    message?: Struct;
    constructor(data?: PartialMessage<SdkMessage>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkMessage";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkMessage;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkMessage;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkMessage;
    static equals(a: SdkMessage | PlainMessage<SdkMessage> | undefined, b: SdkMessage | PlainMessage<SdkMessage> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.AgentMessage
 */
export declare class AgentMessage extends Message<AgentMessage> {
    /**
     * @generated from field: string type = 1;
     */
    type: string;
    /**
     * @generated from field: string uuid = 2;
     */
    uuid: string;
    /**
     * @generated from field: string agent_id = 3;
     */
    agentId: string;
    /**
     * @generated from field: google.protobuf.Struct message = 4;
     */
    message?: Struct;
    constructor(data?: PartialMessage<AgentMessage>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.AgentMessage";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AgentMessage;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AgentMessage;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AgentMessage;
    static equals(a: AgentMessage | PlainMessage<AgentMessage> | undefined, b: AgentMessage | PlainMessage<AgentMessage> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.SdkArtifact
 */
export declare class SdkArtifact extends Message<SdkArtifact> {
    /**
     * @generated from field: string path = 1;
     */
    path: string;
    /**
     * @generated from field: uint64 size_bytes = 2;
     */
    sizeBytes: bigint;
    /**
     * @generated from field: string updated_at = 3;
     */
    updatedAt: string;
    constructor(data?: PartialMessage<SdkArtifact>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.SdkArtifact";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SdkArtifact;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SdkArtifact;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SdkArtifact;
    static equals(a: SdkArtifact | PlainMessage<SdkArtifact> | undefined, b: SdkArtifact | PlainMessage<SdkArtifact> | undefined): boolean;
}
/**
 * @generated from message sdk.v1.DownloadArtifactChunk
 */
export declare class DownloadArtifactChunk extends Message<DownloadArtifactChunk> {
    /**
     * @generated from field: bytes data = 1;
     */
    data: Uint8Array<ArrayBuffer>;
    constructor(data?: PartialMessage<DownloadArtifactChunk>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "sdk.v1.DownloadArtifactChunk";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DownloadArtifactChunk;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DownloadArtifactChunk;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DownloadArtifactChunk;
    static equals(a: DownloadArtifactChunk | PlainMessage<DownloadArtifactChunk> | undefined, b: DownloadArtifactChunk | PlainMessage<DownloadArtifactChunk> | undefined): boolean;
}
//# sourceMappingURL=sdk_messages_pb.d.ts.map