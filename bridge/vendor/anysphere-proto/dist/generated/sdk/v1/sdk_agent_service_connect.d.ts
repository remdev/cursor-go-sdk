import { ArchiveAgentRequest, ArchiveAgentResponse, CancelRunRequest, CancelRunResponse, CloseAgentRequest, CloseAgentResponse, CreateAgentRequest, CreateAgentResponse, DeleteAgentRequest, DeleteAgentResponse, DownloadArtifactRequest, GetAgentRequest, GetAgentResponse, GetRunConversationRequest, GetRunConversationResponse, GetRunRequest, GetRunResponse, ListAgentMessagesRequest, ListAgentMessagesResponse, ListAgentsRequest, ListAgentsResponse, ListArtifactsRequest, ListArtifactsResponse, ListRunsRequest, ListRunsResponse, ObserveRunRequest, ReloadAgentRequest, ReloadAgentResponse, ResumeAgentRequest, ResumeAgentResponse, SendRequest, UnarchiveAgentRequest, UnarchiveAgentResponse, WaitLiveRunRequest, WaitLiveRunResponse } from "./sdk_agent_service_pb.js";
import { MethodKind } from "@bufbuild/protobuf";
import { DownloadArtifactChunk, RunStreamMessage } from "./sdk_messages_pb.js";
/**
 * Agent, run, streaming, cancellation, and artifact APIs.
 *
 * @generated from service sdk.v1.SdkAgentService
 */
export declare const SdkAgentService: {
    readonly typeName: "sdk.v1.SdkAgentService";
    readonly methods: {
        /**
         * @generated from rpc sdk.v1.SdkAgentService.CreateAgent
         */
        readonly createAgent: {
            readonly name: "CreateAgent";
            readonly I: typeof CreateAgentRequest;
            readonly O: typeof CreateAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ResumeAgent
         */
        readonly resumeAgent: {
            readonly name: "ResumeAgent";
            readonly I: typeof ResumeAgentRequest;
            readonly O: typeof ResumeAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ReloadAgent
         */
        readonly reloadAgent: {
            readonly name: "ReloadAgent";
            readonly I: typeof ReloadAgentRequest;
            readonly O: typeof ReloadAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.CloseAgent
         */
        readonly closeAgent: {
            readonly name: "CloseAgent";
            readonly I: typeof CloseAgentRequest;
            readonly O: typeof CloseAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.Send
         */
        readonly send: {
            readonly name: "Send";
            readonly I: typeof SendRequest;
            readonly O: typeof RunStreamMessage;
            readonly kind: MethodKind.ServerStreaming;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.WaitLiveRun
         */
        readonly waitLiveRun: {
            readonly name: "WaitLiveRun";
            readonly I: typeof WaitLiveRunRequest;
            readonly O: typeof WaitLiveRunResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.GetRun
         */
        readonly getRun: {
            readonly name: "GetRun";
            readonly I: typeof GetRunRequest;
            readonly O: typeof GetRunResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ListRuns
         */
        readonly listRuns: {
            readonly name: "ListRuns";
            readonly I: typeof ListRunsRequest;
            readonly O: typeof ListRunsResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.GetRunConversation
         */
        readonly getRunConversation: {
            readonly name: "GetRunConversation";
            readonly I: typeof GetRunConversationRequest;
            readonly O: typeof GetRunConversationResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ObserveRun
         */
        readonly observeRun: {
            readonly name: "ObserveRun";
            readonly I: typeof ObserveRunRequest;
            readonly O: typeof RunStreamMessage;
            readonly kind: MethodKind.ServerStreaming;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.CancelRun
         */
        readonly cancelRun: {
            readonly name: "CancelRun";
            readonly I: typeof CancelRunRequest;
            readonly O: typeof CancelRunResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.GetAgent
         */
        readonly getAgent: {
            readonly name: "GetAgent";
            readonly I: typeof GetAgentRequest;
            readonly O: typeof GetAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ListAgents
         */
        readonly listAgents: {
            readonly name: "ListAgents";
            readonly I: typeof ListAgentsRequest;
            readonly O: typeof ListAgentsResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ArchiveAgent
         */
        readonly archiveAgent: {
            readonly name: "ArchiveAgent";
            readonly I: typeof ArchiveAgentRequest;
            readonly O: typeof ArchiveAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.UnarchiveAgent
         */
        readonly unarchiveAgent: {
            readonly name: "UnarchiveAgent";
            readonly I: typeof UnarchiveAgentRequest;
            readonly O: typeof UnarchiveAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.DeleteAgent
         */
        readonly deleteAgent: {
            readonly name: "DeleteAgent";
            readonly I: typeof DeleteAgentRequest;
            readonly O: typeof DeleteAgentResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ListAgentMessages
         */
        readonly listAgentMessages: {
            readonly name: "ListAgentMessages";
            readonly I: typeof ListAgentMessagesRequest;
            readonly O: typeof ListAgentMessagesResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.ListArtifacts
         */
        readonly listArtifacts: {
            readonly name: "ListArtifacts";
            readonly I: typeof ListArtifactsRequest;
            readonly O: typeof ListArtifactsResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc sdk.v1.SdkAgentService.DownloadArtifact
         */
        readonly downloadArtifact: {
            readonly name: "DownloadArtifact";
            readonly I: typeof DownloadArtifactRequest;
            readonly O: typeof DownloadArtifactChunk;
            readonly kind: MethodKind.ServerStreaming;
        };
    };
};
//# sourceMappingURL=sdk_agent_service_connect.d.ts.map