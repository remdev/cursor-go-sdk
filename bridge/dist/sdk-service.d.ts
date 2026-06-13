import type { ArchiveAgentRequest, CancelRunRequest, CloseAgentRequest, CreateAgentRequest, DeleteAgentRequest, DownloadArtifactRequest, GetAgentRequest, GetRunConversationRequest, GetRunRequest, ListAgentMessagesRequest, ListAgentsRequest, ListArtifactsRequest, ListRunsRequest, ObserveRunRequest, ReloadAgentRequest, ResumeAgentRequest, SendRequest, UnarchiveAgentRequest, WaitLiveRunRequest } from "@anysphere/proto/sdk/v1/sdk_agent_service_pb.js";
import { ArchiveAgentResponse, CancelRunResponse, CloseAgentResponse, CreateAgentResponse, DeleteAgentResponse, GetAgentResponse, GetRunConversationResponse, GetRunResponse, ListAgentMessagesResponse, ListAgentsResponse, ListArtifactsResponse, ListRunsResponse, ReloadAgentResponse, ResumeAgentResponse, UnarchiveAgentResponse, WaitLiveRunResponse } from "@anysphere/proto/sdk/v1/sdk_agent_service_pb.js";
import type { GetVersionRequest, PingRequest, SetToolCallbackRequest, ShutdownRequest } from "@anysphere/proto/sdk/v1/sdk_bridge_control_service_pb.js";
import { GetVersionResponse, PingResponse, SetToolCallbackResponse, ShutdownResponse } from "@anysphere/proto/sdk/v1/sdk_bridge_control_service_pb.js";
import type { ListModelsRequest, ListRepositoriesRequest, MeRequest } from "@anysphere/proto/sdk/v1/sdk_cursor_service_pb.js";
import { ListModelsResponse, ListRepositoriesResponse, MeResponse } from "@anysphere/proto/sdk/v1/sdk_cursor_service_pb.js";
import type { RunStreamMessage } from "@anysphere/proto/sdk/v1/sdk_messages_pb.js";
import { CursorSdkBridgeRegistry } from "./registry.js";
export interface CursorSdkBridgeServiceOptions {
    registry?: CursorSdkBridgeRegistry;
    shutdown?: (graceSeconds: number) => Promise<void> | void;
    /**
     * Absolute workspace this bridge was launched for (`--workspace`). Used as the
     * default `cwd` for local-store operations (create / list / get) when the
     * caller omits one, so reads find agents created here without an explicit
     * per-call `cwd`. The bridge runs as its own subprocess, so `process.cwd()`
     * is the bridge's directory rather than the caller's workspace; this restores
     * the in-process `@cursor/sdk` semantics where `cwd` defaults to the working
     * directory. Per-call `cwd` still overrides.
     */
    workspaceRef?: string;
}
export declare function createSdkAgentService(options?: CursorSdkBridgeServiceOptions): {
    createAgent: (request: CreateAgentRequest) => Promise<CreateAgentResponse>;
    resumeAgent: (request: ResumeAgentRequest) => Promise<ResumeAgentResponse>;
    reloadAgent: (request: ReloadAgentRequest) => Promise<ReloadAgentResponse>;
    closeAgent: (request: CloseAgentRequest) => Promise<CloseAgentResponse>;
    send: (request: SendRequest) => AsyncGenerator<RunStreamMessage, void, any>;
    waitLiveRun: (request: WaitLiveRunRequest) => Promise<WaitLiveRunResponse>;
    getRun: (request: GetRunRequest) => Promise<GetRunResponse>;
    listRuns: (request: ListRunsRequest) => Promise<ListRunsResponse>;
    getRunConversation: (request: GetRunConversationRequest) => Promise<GetRunConversationResponse>;
    observeRun: (request: ObserveRunRequest) => AsyncGenerator<RunStreamMessage, void, any>;
    cancelRun: (request: CancelRunRequest) => Promise<CancelRunResponse>;
    getAgent: (request: GetAgentRequest) => Promise<GetAgentResponse>;
    listAgents: (request: ListAgentsRequest) => Promise<ListAgentsResponse>;
    archiveAgent: (request: ArchiveAgentRequest) => Promise<ArchiveAgentResponse>;
    unarchiveAgent: (request: UnarchiveAgentRequest) => Promise<UnarchiveAgentResponse>;
    deleteAgent: (request: DeleteAgentRequest) => Promise<DeleteAgentResponse>;
    listAgentMessages: (request: ListAgentMessagesRequest) => Promise<ListAgentMessagesResponse>;
    listArtifacts: (request: ListArtifactsRequest) => Promise<ListArtifactsResponse>;
    downloadArtifact: (request: DownloadArtifactRequest) => AsyncGenerator<import("@anysphere/proto/sdk/v1/sdk_messages_pb.js").DownloadArtifactChunk, void, unknown>;
};
export declare function createSdkCursorService(): {
    me: (request: MeRequest) => Promise<MeResponse>;
    listModels: (request: ListModelsRequest) => Promise<ListModelsResponse>;
    listRepositories: (request: ListRepositoriesRequest) => Promise<ListRepositoriesResponse>;
};
export declare function createSdkBridgeControlService(options?: CursorSdkBridgeServiceOptions): {
    ping: (_request: PingRequest) => Promise<PingResponse>;
    shutdown: (request: ShutdownRequest) => Promise<ShutdownResponse>;
    getVersion: (_request: GetVersionRequest) => Promise<GetVersionResponse>;
    setToolCallback: (request: SetToolCallbackRequest) => Promise<SetToolCallbackResponse>;
};
//# sourceMappingURL=sdk-service.d.ts.map