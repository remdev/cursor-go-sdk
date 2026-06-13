import { randomUUID } from "node:crypto";
import {
    ArchiveAgentResponse,
    CancelRunResponse,
    CloseAgentResponse,
    CreateAgentResponse,
    DeleteAgentResponse,
    GetAgentResponse,
    GetRunConversationResponse,
    GetRunResponse,
    ListAgentMessagesResponse,
    ListAgentsResponse,
    ListArtifactsResponse,
    ListRunsResponse,
    ReloadAgentResponse,
    ResumeAgentResponse,
    UnarchiveAgentResponse,
    WaitLiveRunResponse,
    type CreateAgentRequest,
    type SendRequest,
} from "../gen/ts/sdk/v1/sdk_agent_service_pb.js";
import {
    GetVersionResponse,
    PingResponse,
    SetToolCallbackResponse,
    ShutdownResponse,
} from "../gen/ts/sdk/v1/sdk_bridge_control_service_pb.js";
import {
    ListModelsResponse,
    ListRepositoriesResponse,
    MeResponse,
} from "../gen/ts/sdk/v1/sdk_cursor_service_pb.js";
import {
    Runtime,
    type AgentOptions,
} from "../gen/ts/sdk/v1/sdk_messages_pb.js";
import { Code, ConnectError } from "@connectrpc/connect";
import { Agent, Cursor, type Run } from "@cursor/sdk";
import {
    assertCustomToolsConfigured,
    assertCustomToolsLocalOnly,
    hasCustomToolDefinitions,
    protoCustomToolDefinitionsToHost,
    resetBridgeToolCallbackClient,
} from "./bridge-custom-tools.js";
import {
    CURSOR_SDK_BRIDGE_CAPABILITIES,
    CURSOR_SDK_BRIDGE_PROTOCOL_VERSION,
    CURSOR_SDK_BRIDGE_VERSION,
} from "./constants.js";
import { CursorSdkBridgeRegistry } from "./registry.js";
import {
    protoAgentMessagesOptionsToSdk,
    protoAgentOperationOptionsToSdk,
    protoAgentOptionsToSdk,
    protoCursorRequestOptionsToSdk,
    protoGetRunOptionsToSdk,
    protoListAgentsOptionsToSdk,
    protoListRunsOptionsToSdk,
    protoSendOptionsToSdk,
    protoUserMessageToSdk,
    sdkAgentInfoToProto,
    sdkAgentMessagesToProto,
    sdkAgentToCreateResponse,
    sdkArtifactBufferToChunks,
    sdkArtifactsToProto,
    sdkConversationStepToRunStreamMessage,
    sdkInteractionUpdateToRunStreamMessage,
    sdkListAgentsResponseItemsToProto,
    sdkMessageToRunStreamMessage,
    sdkModelsToProto,
    sdkRepositoriesToProto,
    sdkRunDoneToRunStreamMessage,
    sdkRunResultToProto,
    sdkRunResultToRunStreamMessage,
    sdkRunToSnapshot,
    sdkUserToProto,
} from "./sdk-converters.js";
import { setBridgeToolCallbackConfig } from "./tool-callback-config.js";
import type { CursorSdkBridgeServiceOptions, SdkBridgeControlServiceOptions } from "./types.js";
import type { RunStreamMessage } from "../gen/ts/sdk/v1/sdk_messages_pb.js";

const ARTIFACT_CHUNK_BYTES = 1024 * 1024;

function isCloudAgentId(agentId: string): boolean {
    return agentId.startsWith("bc-");
}

export function createSdkAgentService(options: CursorSdkBridgeServiceOptions = {}) {
    const registry = options.registry ?? new CursorSdkBridgeRegistry();
    const defaultCwd =
        options.workspaceRef && options.workspaceRef.length > 0
            ? options.workspaceRef
            : undefined;

    const getAgent = (agentId: string) => {
        const agent = registry.getAgent(agentId);
        if (!agent) {
            throw new ConnectError(`Unknown agent: ${agentId}`, Code.NotFound);
        }
        return agent;
    };

    const getRun = (runId: string) => {
        const run = registry.getRun(runId);
        if (!run) {
            throw new ConnectError(`Unknown run: ${runId}`, Code.NotFound);
        }
        return run;
    };

    return {
        createAgent: async (request: CreateAgentRequest) => {
            if (request.idempotencyKey && !request.options?.cloud) {
                throw new ConnectError(
                    "Idempotency-Key is only supported for cloud CreateAgent in v1",
                    Code.Unimplemented,
                );
            }
            validateCreateAgentRequest(request);
            validateAgentCustomTools(request.options);
            const hostCustomTools = protoCustomToolDefinitionsToHost(
                request.options?.local?.customTools,
            );
            const customToolsAgentId = hostCustomTools
                ? request.options?.agentId || randomUUID()
                : undefined;
            const agentOptions = protoAgentOptionsToSdk(
                request.options,
                defaultCwd,
                customToolsAgentId,
            );
            if (customToolsAgentId && !request.options?.agentId) {
                agentOptions.agentId = customToolsAgentId;
            }
            const agent = await Agent.create({
                ...agentOptions,
                ...(request.idempotencyKey ? { idempotencyKey: request.idempotencyKey } : {}),
            });
            await registry.registerAgent(agent, { cloud: Boolean(request.options?.cloud) });
            return new CreateAgentResponse(sdkAgentToCreateResponse(agent));
        },

        resumeAgent: async (request) => {
            validateAgentCustomTools(request.options);
            const customToolsAgentId = hasCustomToolDefinitions(request.options?.local?.customTools)
                ? request.agentId
                : undefined;
            const agentOptions = protoAgentOptionsToSdk(
                request.options,
                defaultCwd,
                customToolsAgentId,
            );
            const agent = await Agent.resume(request.agentId, agentOptions);
            await registry.registerAgent(agent, {
                cloud: Boolean(agentOptions.cloud) || isCloudAgentId(request.agentId),
            });
            return new ResumeAgentResponse(sdkAgentToCreateResponse(agent));
        },

        reloadAgent: async (request) => {
            await getAgent(request.agentId).reload();
            return new ReloadAgentResponse();
        },

        closeAgent: async (request) => {
            void getAgent(request.agentId);
            await registry.disposeAgent(request.agentId);
            return new CloseAgentResponse();
        },

        send: async function* (request: SendRequest) {
            const agent = getAgent(request.agentId);
            const isCloudAgent = registry.isCloudAgent(request.agentId);
            if (request.idempotencyKey && !isCloudAgent) {
                throw new ConnectError(
                    "Idempotency-Key is only supported for cloud Send in v1",
                    Code.Unimplemented,
                );
            }
            validateSendRequest(request, isCloudAgent);
            const wantsMergedStream = Boolean(
                request.options?.enableDeltas || request.options?.enableSteps,
            );
            const streamQueue = wantsMergedStream ? new RunStreamQueue() : undefined;
            const run = await agent.send(protoUserMessageToSdk(request.message), {
                ...protoSendOptionsToSdk(request.options),
                ...(request.options?.enableDeltas && streamQueue
                    ? {
                        onDelta: ({ update }) => {
                            streamQueue.enqueue((offset) =>
                                sdkInteractionUpdateToRunStreamMessage(update, offset),
                            );
                        },
                    }
                    : {}),
                ...(request.options?.enableSteps && streamQueue
                    ? {
                        onStep: ({ step }) => {
                            streamQueue.enqueue((offset) =>
                                sdkConversationStepToRunStreamMessage(step, offset),
                            );
                        },
                    }
                    : {}),
                ...(request.idempotencyKey ? { idempotencyKey: request.idempotencyKey } : {}),
            });
            if (agent.agentId !== request.agentId) {
                await registry.registerAgent(agent, { cloud: isCloudAgent });
            }
            registry.registerRun(run);
            if (streamQueue) {
                yield* yieldMergedRunStream(run, streamQueue);
            } else {
                yield* yieldRunStream(run);
            }
        },

        waitLiveRun: async (request) => {
            const run = getRun(request.runId);
            return new WaitLiveRunResponse({
                result: sdkRunResultToProto(await run.wait(), run.agentId, run.createdAt),
            });
        },

        getRun: async (request) => {
            if (
                request.options?.runtime === Runtime.CLOUD &&
                request.options.agentId.length === 0
            ) {
                throw new ConnectError("Cloud getRun requests require agentId", Code.InvalidArgument);
            }
            const run = await Agent.getRun(
                request.runId,
                protoGetRunOptionsToSdk(request.options, defaultCwd),
            );
            registry.registerRun(run);
            return new GetRunResponse({ run: sdkRunToSnapshot(run) });
        },

        listRuns: async (request) => {
            const response = await Agent.listRuns(
                request.agentId,
                protoListRunsOptionsToSdk(request.options, defaultCwd),
            );
            registry.registerRuns(response.items);
            return new ListRunsResponse({
                items: response.items.map(sdkRunToSnapshot),
                nextCursor: response.nextCursor ?? "",
            });
        },

        getRunConversation: async (request) =>
            new GetRunConversationResponse({
                conversationJson: JSON.stringify(await getRun(request.runId).conversation()),
            }),

        observeRun: async function* (request) {
            let run = registry.getRun(request.runId);
            if (!run) {
                run = await Agent.getRun(request.runId);
                registry.registerRun(run);
            }
            yield* yieldRunStream(run, request.afterOffset);
        },

        cancelRun: async (request) => {
            const liveRun = registry.getRun(request.runId);
            if (liveRun) {
                await liveRun.cancel();
            } else if (Agent.cancelRun) {
                await Agent.cancelRun(
                    request.runId,
                    request.agentId
                        ? { runtime: "cloud", agentId: request.agentId }
                        : { runtime: "local" },
                );
            } else {
                throw new ConnectError(
                    "Detached run cancellation requires @cursor/sdk Agent.cancelRun support",
                    Code.Unimplemented,
                );
            }
            return new CancelRunResponse();
        },

        getAgent: async (request) =>
            new GetAgentResponse({
                agent: sdkAgentInfoToProto(
                    await Agent.get(
                        request.agentId,
                        protoAgentOperationOptionsToSdk(request.options, defaultCwd),
                    ),
                ),
            }),

        listAgents: async (request) => {
            const response = await Agent.list(
                protoListAgentsOptionsToSdk(request.options, defaultCwd),
            );
            return new ListAgentsResponse({
                items: sdkListAgentsResponseItemsToProto(response.items),
                nextCursor: response.nextCursor ?? "",
            });
        },

        archiveAgent: async (request) => {
            await Agent.archive(
                request.agentId,
                protoAgentOperationOptionsToSdk(request.options, defaultCwd),
            );
            return new ArchiveAgentResponse();
        },

        unarchiveAgent: async (request) => {
            await Agent.unarchive(
                request.agentId,
                protoAgentOperationOptionsToSdk(request.options, defaultCwd),
            );
            return new UnarchiveAgentResponse();
        },

        deleteAgent: async (request) => {
            await Agent.delete(
                request.agentId,
                protoAgentOperationOptionsToSdk(request.options, defaultCwd),
            );
            return new DeleteAgentResponse();
        },

        listAgentMessages: async (request) => {
            if (
                request.options?.runtime === Runtime.CLOUD ||
                (request.options?.runtime !== Runtime.LOCAL && isCloudAgentId(request.agentId))
            ) {
                throw new ConnectError(
                    "Cloud ListAgentMessages is not supported by @cursor/sdk yet",
                    Code.Unimplemented,
                );
            }
            return new ListAgentMessagesResponse({
                messages: sdkAgentMessagesToProto(
                    await Agent.messages.list(
                        request.agentId,
                        protoAgentMessagesOptionsToSdk(request.options, defaultCwd),
                    ),
                ),
            });
        },

        listArtifacts: async (request) =>
            new ListArtifactsResponse({
                artifacts: sdkArtifactsToProto(await getAgent(request.agentId).listArtifacts()),
            }),

        downloadArtifact: async function* (request) {
            const buffer = await getAgent(request.agentId).downloadArtifact(request.path);
            for (const chunk of sdkArtifactBufferToChunks(buffer, ARTIFACT_CHUNK_BYTES)) {
                yield chunk;
            }
        },
    };
}

async function* yieldRunStream(run: Run, afterOffset?: string): AsyncGenerator<RunStreamMessage> {
    const startAfter = parseOffset(afterOffset);
    let offset = 0;
    for await (const message of run.stream()) {
        offset += 1;
        if (offset <= startAfter) {
            continue;
        }
        yield sdkMessageToRunStreamMessage(message, String(offset));
    }
    const result = await run.wait();
    offset += 1;
    if (offset > startAfter) {
        yield sdkRunResultToRunStreamMessage(run, result, String(offset));
    }
    offset += 1;
    if (offset > startAfter) {
        yield sdkRunDoneToRunStreamMessage(run, String(offset));
    }
}

async function* yieldMergedRunStream(
    run: Run,
    streamQueue: RunStreamQueue,
    afterOffset?: string,
): AsyncGenerator<RunStreamMessage> {
    const stream = run.stream();
    const producer = (async () => {
        try {
            for await (const message of stream) {
                if (streamQueue.isClosed) {
                    return;
                }
                streamQueue.enqueue((offset) => sdkMessageToRunStreamMessage(message, offset));
                await streamQueue.waitForDrain();
            }
            if (streamQueue.isClosed) {
                return;
            }
            const result = await run.wait();
            streamQueue.enqueue((offset) => sdkRunResultToRunStreamMessage(run, result, offset));
            streamQueue.enqueue((offset) => sdkRunDoneToRunStreamMessage(run, offset));
        } catch (error) {
            streamQueue.fail(error);
            return;
        }
        streamQueue.close();
    })();
    const startAfter = parseOffset(afterOffset);
    try {
        for await (const event of streamQueue) {
            if (parseOffset(event.offset) > startAfter) {
                yield event;
            }
        }
    } finally {
        streamQueue.close();
        await stream.return?.();
        await producer;
    }
}

const RUN_STREAM_QUEUE_HIGH_WATER_MARK = 1024;
const RUN_STREAM_QUEUE_HARD_LIMIT = 4096;

class RunStreamQueue implements AsyncIterable<RunStreamMessage> {
    private pending: RunStreamMessage[] = [];
    private waiters: Array<() => void> = [];
    private drainWaiters: Array<() => void> = [];
    private offset = 0;
    private closed = false;
    private hasFailed = false;
    private failure: unknown;

    enqueue(factory: (offset: string) => RunStreamMessage): void {
        if (this.closed) {
            return;
        }
        if (this.pending.length >= RUN_STREAM_QUEUE_HARD_LIMIT) {
            this.fail(
                new ConnectError(
                    "Merged run stream consumer fell behind; the bridge dropped the connection to bound memory",
                    Code.ResourceExhausted,
                ),
            );
            return;
        }
        this.offset += 1;
        this.pending.push(factory(String(this.offset)));
        this.notify();
    }

    async waitForDrain(): Promise<void> {
        if (this.closed || this.pending.length < RUN_STREAM_QUEUE_HIGH_WATER_MARK) {
            return;
        }
        await new Promise<void>((resolve) => this.drainWaiters.push(resolve));
    }

    close(): void {
        if (this.closed) {
            return;
        }
        this.closed = true;
        this.notify();
        this.notifyDrain();
    }

    fail(error: unknown): void {
        this.hasFailed = true;
        this.failure = error;
        this.closed = true;
        this.notify();
        this.notifyDrain();
    }

    get isClosed(): boolean {
        return this.closed;
    }

    async *[Symbol.asyncIterator](): AsyncGenerator<RunStreamMessage> {
        while (true) {
            const next = this.pending.shift();
            if (next) {
                if (this.pending.length < RUN_STREAM_QUEUE_HIGH_WATER_MARK) {
                    this.notifyDrain();
                }
                yield next;
                continue;
            }
            if (this.hasFailed) {
                throw this.failure;
            }
            if (this.closed) {
                return;
            }
            await new Promise<void>((resolve) => this.waiters.push(resolve));
        }
    }

    private notify(): void {
        for (const waiter of this.waiters.splice(0)) {
            waiter();
        }
    }

    private notifyDrain(): void {
        for (const waiter of this.drainWaiters.splice(0)) {
            waiter();
        }
    }
}

function validateSendRequest(request: SendRequest, isCloudAgent: boolean): void {
    if (isCloudAgent) {
        return;
    }
    for (const image of request.message?.images ?? []) {
        if (image.source.case === "url") {
            throw new ConnectError(
                "URL images are only supported for cloud-routed agents",
                Code.InvalidArgument,
            );
        }
    }
}

function validateAgentCustomTools(options: AgentOptions | undefined): void {
    assertCustomToolsLocalOnly({
        cloud: Boolean(options?.cloud),
        localCustomTools: options?.local?.customTools,
    });
    assertCustomToolsConfigured(options?.local?.customTools);
}

function validateCreateAgentRequest(request: CreateAgentRequest): void {
    const envVars = request.options?.cloud?.envVars ?? {};
    if (Object.keys(envVars).length === 0) {
        return;
    }
    if (request.options?.agentId) {
        throw new ConnectError(
            "cloud.envVars cannot be used with a caller-supplied agentId",
            Code.InvalidArgument,
        );
    }
    for (const name of Object.keys(envVars)) {
        if (name.startsWith("CURSOR_")) {
            throw new ConnectError(
                "cloud.envVars names cannot start with CURSOR_",
                Code.InvalidArgument,
            );
        }
    }
}

function parseOffset(offset: string | undefined): number {
    if (!offset) {
        return 0;
    }
    const parsed = Number.parseInt(offset, 10);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : 0;
}

export function createSdkCursorService() {
    return {
        me: async (request) => {
            const options = requireCursorRequestApiKey(request.options);
            return new MeResponse({
                user: sdkUserToProto(await Cursor.me(options)),
            });
        },
        listModels: async (request) => {
            const options = requireCursorRequestApiKey(request.options);
            return new ListModelsResponse({
                items: sdkModelsToProto(await Cursor.models.list(options)),
            });
        },
        listRepositories: async (request) => {
            const options = requireCursorRequestApiKey(request.options);
            return new ListRepositoriesResponse({
                items: sdkRepositoriesToProto(await Cursor.repositories.list(options)),
            });
        },
    };
}

function requireCursorRequestApiKey(options: Parameters<typeof protoCursorRequestOptionsToSdk>[0]) {
    const converted = protoCursorRequestOptionsToSdk(options);
    if (converted === undefined || !converted.apiKey) {
        throw new ConnectError(
            "API key is required for cloud catalog calls.",
            Code.Unauthenticated,
        );
    }
    return converted;
}

export function createSdkBridgeControlService(options: SdkBridgeControlServiceOptions = {}) {
    return {
        ping: async (_request) => new PingResponse({ message: "pong" }),
        shutdown: async (request) => {
            await options.shutdown?.(request.graceSeconds);
            return new ShutdownResponse();
        },
        getVersion: async (_request) =>
            new GetVersionResponse({
                bridgeVersion: CURSOR_SDK_BRIDGE_VERSION,
                protocolVersion: CURSOR_SDK_BRIDGE_PROTOCOL_VERSION,
                capabilities: CURSOR_SDK_BRIDGE_CAPABILITIES,
            }),
        setToolCallback: async (request) => {
            const url = request.url.trim();
            const authToken = request.authToken.trim();
            if (!url) {
                setBridgeToolCallbackConfig(undefined);
                resetBridgeToolCallbackClient();
                return new SetToolCallbackResponse();
            }
            if (!authToken) {
                throw new ConnectError(
                    "SetToolCallback requires an auth_token when url is set.",
                    Code.InvalidArgument,
                );
            }
            setBridgeToolCallbackConfig({ url, authToken });
            resetBridgeToolCallbackClient();
            return new SetToolCallbackResponse();
        },
    };
}
