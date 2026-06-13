import { randomUUID } from "node:crypto";
import { ArchiveAgentResponse, CancelRunResponse, CloseAgentResponse, CreateAgentResponse, DeleteAgentResponse, GetAgentResponse, GetRunConversationResponse, GetRunResponse, ListAgentMessagesResponse, ListAgentsResponse, ListArtifactsResponse, ListRunsResponse, ReloadAgentResponse, ResumeAgentResponse, UnarchiveAgentResponse, WaitLiveRunResponse, } from "@anysphere/proto/sdk/v1/sdk_agent_service_pb.js";
import { GetVersionResponse, PingResponse, SetToolCallbackResponse, ShutdownResponse, } from "@anysphere/proto/sdk/v1/sdk_bridge_control_service_pb.js";
import { ListModelsResponse, ListRepositoriesResponse, MeResponse, } from "@anysphere/proto/sdk/v1/sdk_cursor_service_pb.js";
import { Runtime } from "@anysphere/proto/sdk/v1/sdk_messages_pb.js";
import { Code, ConnectError } from "@connectrpc/connect";
import { Agent, Cursor } from "@cursor/sdk";
import { assertCustomToolsConfigured, assertCustomToolsLocalOnly, hasCustomToolDefinitions, protoCustomToolDefinitionsToHost, resetBridgeToolCallbackClient, } from "./bridge-custom-tools.js";
import { CURSOR_SDK_BRIDGE_CAPABILITIES, CURSOR_SDK_BRIDGE_PROTOCOL_VERSION, CURSOR_SDK_BRIDGE_VERSION, } from "./constants.js";
import { CursorSdkBridgeRegistry } from "./registry.js";
import { protoAgentMessagesOptionsToSdk, protoAgentOperationOptionsToSdk, protoAgentOptionsToSdk, protoCursorRequestOptionsToSdk, protoGetRunOptionsToSdk, protoListAgentsOptionsToSdk, protoListRunsOptionsToSdk, protoSendOptionsToSdk, protoUserMessageToSdk, sdkAgentInfoToProto, sdkAgentMessagesToProto, sdkAgentToCreateResponse, sdkArtifactBufferToChunks, sdkArtifactsToProto, sdkConversationStepToRunStreamMessage, sdkInteractionUpdateToRunStreamMessage, sdkListAgentsResponseItemsToProto, sdkMessageToRunStreamMessage, sdkModelsToProto, sdkRepositoriesToProto, sdkRunDoneToRunStreamMessage, sdkRunResultToProto, sdkRunResultToRunStreamMessage, sdkRunToSnapshot, sdkUserToProto, } from "./sdk-converters.js";
import { setBridgeToolCallbackConfig } from "./tool-callback-config.js";
const ARTIFACT_CHUNK_BYTES = 1024 * 1024;
const BridgeAgentApi = Agent;
const BridgeAgentMessages = Agent.messages;
function isCloudAgentId(agentId) {
    return agentId.startsWith("bc-");
}
export function createSdkAgentService(options = {}) {
    var _a;
    const registry = (_a = options.registry) !== null && _a !== void 0 ? _a : new CursorSdkBridgeRegistry();
    // Default `cwd` for local-store ops when the caller omits one. See
    // `CursorSdkBridgeServiceOptions.workspaceRef`.
    const defaultCwd = options.workspaceRef && options.workspaceRef.length > 0
        ? options.workspaceRef
        : undefined;
    const getAgent = (agentId) => {
        const agent = registry.getAgent(agentId);
        if (!agent) {
            throw new ConnectError(`Unknown agent: ${agentId}`, Code.NotFound);
        }
        return agent;
    };
    const getRun = (runId) => {
        const run = registry.getRun(runId);
        if (!run) {
            throw new ConnectError(`Unknown run: ${runId}`, Code.NotFound);
        }
        return run;
    };
    return {
        createAgent: async (request) => {
            var _a, _b, _c, _d, _e, _f;
            if (request.idempotencyKey && !((_a = request.options) === null || _a === void 0 ? void 0 : _a.cloud)) {
                throw new ConnectError("Idempotency-Key is only supported for cloud CreateAgent in v1", Code.Unimplemented);
            }
            validateCreateAgentRequest(request);
            validateAgentCustomTools(request.options);
            const hostCustomTools = protoCustomToolDefinitionsToHost((_c = (_b = request.options) === null || _b === void 0 ? void 0 : _b.local) === null || _c === void 0 ? void 0 : _c.customTools);
            const customToolsAgentId = hostCustomTools
                ? ((_d = request.options) === null || _d === void 0 ? void 0 : _d.agentId) || randomUUID()
                : undefined;
            const agentOptions = protoAgentOptionsToSdk(request.options, defaultCwd, customToolsAgentId);
            if (customToolsAgentId && !((_e = request.options) === null || _e === void 0 ? void 0 : _e.agentId)) {
                agentOptions.agentId = customToolsAgentId;
            }
            const agent = await BridgeAgentApi.create(Object.assign(Object.assign({}, agentOptions), (request.idempotencyKey
                ? { idempotencyKey: request.idempotencyKey }
                : {})));
            await registry.registerAgent(agent, {
                cloud: Boolean((_f = request.options) === null || _f === void 0 ? void 0 : _f.cloud),
            });
            return new CreateAgentResponse(sdkAgentToCreateResponse(agent));
        },
        resumeAgent: async (request) => {
            var _a, _b;
            validateAgentCustomTools(request.options);
            const customToolsAgentId = hasCustomToolDefinitions((_b = (_a = request.options) === null || _a === void 0 ? void 0 : _a.local) === null || _b === void 0 ? void 0 : _b.customTools)
                ? request.agentId
                : undefined;
            const agentOptions = protoAgentOptionsToSdk(request.options, defaultCwd, customToolsAgentId);
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
        send: async function* (request) {
            var _a, _b, _c, _d;
            const agent = getAgent(request.agentId);
            const isCloudAgent = registry.isCloudAgent(request.agentId);
            if (request.idempotencyKey && !isCloudAgent) {
                throw new ConnectError("Idempotency-Key is only supported for cloud Send in v1", Code.Unimplemented);
            }
            validateSendRequest(request, isCloudAgent);
            // The merged buffered path is only needed for delta/step callbacks; the
            // default path uses yieldRunStream so for await + yield apply natural
            // backpressure to the SDK source.
            const wantsMergedStream = Boolean(((_a = request.options) === null || _a === void 0 ? void 0 : _a.enableDeltas) || ((_b = request.options) === null || _b === void 0 ? void 0 : _b.enableSteps));
            const streamQueue = wantsMergedStream ? new RunStreamQueue() : undefined;
            const run = await agent.send(protoUserMessageToSdk(request.message), Object.assign(Object.assign(Object.assign(Object.assign({}, protoSendOptionsToSdk(request.options)), (((_c = request.options) === null || _c === void 0 ? void 0 : _c.enableDeltas) && streamQueue
                ? {
                    onDelta: ({ update }) => {
                        streamQueue.enqueue(offset => sdkInteractionUpdateToRunStreamMessage(update, offset));
                    },
                }
                : {})), (((_d = request.options) === null || _d === void 0 ? void 0 : _d.enableSteps) && streamQueue
                ? {
                    onStep: ({ step }) => {
                        streamQueue.enqueue(offset => sdkConversationStepToRunStreamMessage(step, offset));
                    },
                }
                : {})), (request.idempotencyKey
                ? { idempotencyKey: request.idempotencyKey }
                : {})));
            if (agent.agentId !== request.agentId) {
                await registry.registerAgent(agent, { cloud: isCloudAgent });
            }
            registry.registerRun(run);
            if (streamQueue) {
                yield* yieldMergedRunStream(run, streamQueue);
            }
            else {
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
            var _a;
            if (((_a = request.options) === null || _a === void 0 ? void 0 : _a.runtime) === Runtime.CLOUD &&
                request.options.agentId.length === 0) {
                throw new ConnectError("Cloud getRun requests require agentId", Code.InvalidArgument);
            }
            const run = await Agent.getRun(request.runId, protoGetRunOptionsToSdk(request.options, defaultCwd));
            registry.registerRun(run);
            return new GetRunResponse({ run: sdkRunToSnapshot(run) });
        },
        listRuns: async (request) => {
            var _a;
            const response = await Agent.listRuns(request.agentId, protoListRunsOptionsToSdk(request.options, defaultCwd));
            const runs = response.items;
            registry.registerRuns(runs);
            return new ListRunsResponse({
                items: runs.map(sdkRunToSnapshot),
                nextCursor: (_a = response.nextCursor) !== null && _a !== void 0 ? _a : "",
            });
        },
        getRunConversation: async (request) => new GetRunConversationResponse({
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
            }
            else {
                const cancelRun = BridgeAgentApi.cancelRun;
                if (!cancelRun) {
                    throw new ConnectError("Detached run cancellation requires @cursor/sdk Agent.cancelRun support", Code.Unimplemented);
                }
                await cancelRun(request.runId, request.agentId
                    ? { runtime: "cloud", agentId: request.agentId }
                    : { runtime: "local" });
            }
            return new CancelRunResponse();
        },
        getAgent: async (request) => new GetAgentResponse({
            agent: sdkAgentInfoToProto(await Agent.get(request.agentId, protoAgentOperationOptionsToSdk(request.options, defaultCwd))),
        }),
        listAgents: async (request) => {
            var _a;
            const response = await Agent.list(protoListAgentsOptionsToSdk(request.options, defaultCwd));
            return new ListAgentsResponse({
                items: sdkListAgentsResponseItemsToProto(response.items),
                nextCursor: (_a = response.nextCursor) !== null && _a !== void 0 ? _a : "",
            });
        },
        archiveAgent: async (request) => {
            await Agent.archive(request.agentId, protoAgentOperationOptionsToSdk(request.options, defaultCwd));
            return new ArchiveAgentResponse();
        },
        unarchiveAgent: async (request) => {
            await Agent.unarchive(request.agentId, protoAgentOperationOptionsToSdk(request.options, defaultCwd));
            return new UnarchiveAgentResponse();
        },
        deleteAgent: async (request) => {
            await Agent.delete(request.agentId, protoAgentOperationOptionsToSdk(request.options, defaultCwd));
            return new DeleteAgentResponse();
        },
        listAgentMessages: async (request) => {
            var _a, _b;
            if (((_a = request.options) === null || _a === void 0 ? void 0 : _a.runtime) === Runtime.CLOUD ||
                (((_b = request.options) === null || _b === void 0 ? void 0 : _b.runtime) !== Runtime.LOCAL &&
                    isCloudAgentId(request.agentId))) {
                throw new ConnectError("Cloud ListAgentMessages is not supported by @cursor/sdk yet", Code.Unimplemented);
            }
            return new ListAgentMessagesResponse({
                messages: sdkAgentMessagesToProto(await BridgeAgentMessages.list(request.agentId, protoAgentMessagesOptionsToSdk(request.options, defaultCwd))),
            });
        },
        listArtifacts: async (request) => new ListArtifactsResponse({
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
async function* yieldRunStream(run, afterOffset) {
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
async function* yieldMergedRunStream(run, streamQueue, afterOffset) {
    var _a;
    const stream = run.stream();
    const producer = (async () => {
        try {
            for await (const message of stream) {
                if (streamQueue.isClosed) {
                    return;
                }
                streamQueue.enqueue(offset => sdkMessageToRunStreamMessage(message, offset));
                // Pause the SDK source while the consumer drains; this also throttles
                // onDelta / onStep callbacks, which fire synchronously in this tick.
                await streamQueue.waitForDrain();
            }
            // The for-await above can exit via stream.return() when the consumer
            // disconnects (see the cleanup `finally` below). Don't block on
            // run.wait() in that case -- the run may keep going indefinitely and
            // the RPC handler is already tearing down.
            if (streamQueue.isClosed) {
                return;
            }
            const result = await run.wait();
            streamQueue.enqueue(offset => sdkRunResultToRunStreamMessage(run, result, offset));
            streamQueue.enqueue(offset => sdkRunDoneToRunStreamMessage(run, offset));
        }
        catch (error) {
            streamQueue.fail(error);
            return;
        }
        streamQueue.close();
    })();
    const startAfter = parseOffset(afterOffset);
    try {
        for await (const event of streamQueue) {
            const numericOffset = parseOffset(event.offset);
            if (numericOffset > startAfter) {
                yield event;
            }
        }
    }
    finally {
        // If the consumer disconnects mid-stream, close the queue and ask the
        // underlying SDK stream generator to return so its `for await` terminates
        // even when the run is idle (no further messages to emit). Without
        // calling stream.return() the producer's `for await` would suspend on
        // the run's pending promise forever, blocking RPC-handler cleanup.
        streamQueue.close();
        await ((_a = stream.return) === null || _a === void 0 ? void 0 : _a.call(stream));
        await producer;
    }
}
// Producer awaits at the high-water mark; hard cap fails fast for synchronous
// callbacks that can't await. Sized to absorb a few seconds of token-rate
// `enable_deltas=true` activity before backpressure kicks in.
const RUN_STREAM_QUEUE_HIGH_WATER_MARK = 1024;
const RUN_STREAM_QUEUE_HARD_LIMIT = 4096;
class RunStreamQueue {
    pending = [];
    waiters = [];
    drainWaiters = [];
    offset = 0;
    closed = false;
    hasFailed = false;
    failure;
    enqueue(factory) {
        if (this.closed) {
            return;
        }
        if (this.pending.length >= RUN_STREAM_QUEUE_HARD_LIMIT) {
            // Fail closed when synchronous callbacks burst past the hard cap so the
            // bridge process bounds memory rather than the gRPC client crashing it.
            this.fail(new ConnectError("Merged run stream consumer fell behind; the bridge dropped the connection to bound memory", Code.ResourceExhausted));
            return;
        }
        this.offset += 1;
        this.pending.push(factory(String(this.offset)));
        this.notify();
    }
    async waitForDrain() {
        if (this.closed || this.pending.length < RUN_STREAM_QUEUE_HIGH_WATER_MARK) {
            return;
        }
        await new Promise(resolve => this.drainWaiters.push(resolve));
    }
    close() {
        if (this.closed) {
            return;
        }
        this.closed = true;
        this.notify();
        this.notifyDrain();
    }
    fail(error) {
        this.hasFailed = true;
        this.failure = error;
        this.closed = true;
        this.notify();
        this.notifyDrain();
    }
    get isClosed() {
        return this.closed;
    }
    async *[Symbol.asyncIterator]() {
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
            await new Promise(resolve => this.waiters.push(resolve));
        }
    }
    notify() {
        const waiters = this.waiters.splice(0, this.waiters.length);
        for (const waiter of waiters) {
            waiter();
        }
    }
    notifyDrain() {
        const waiters = this.drainWaiters.splice(0, this.drainWaiters.length);
        for (const waiter of waiters) {
            waiter();
        }
    }
}
function validateSendRequest(request, isCloudAgent) {
    var _a;
    var _b;
    if (isCloudAgent) {
        return;
    }
    for (const image of (_b = (_a = request.message) === null || _a === void 0 ? void 0 : _a.images) !== null && _b !== void 0 ? _b : []) {
        if (image.source.case === "url") {
            throw new ConnectError("URL images are only supported for cloud-routed agents", Code.InvalidArgument);
        }
    }
}
function validateAgentCustomTools(options) {
    var _a, _b;
    assertCustomToolsLocalOnly({
        cloud: options === null || options === void 0 ? void 0 : options.cloud,
        localCustomTools: (_a = options === null || options === void 0 ? void 0 : options.local) === null || _a === void 0 ? void 0 : _a.customTools,
    });
    assertCustomToolsConfigured((_b = options === null || options === void 0 ? void 0 : options.local) === null || _b === void 0 ? void 0 : _b.customTools);
}
function validateCreateAgentRequest(request) {
    var _a, _b, _c;
    var _d;
    const envVars = (_d = (_b = (_a = request.options) === null || _a === void 0 ? void 0 : _a.cloud) === null || _b === void 0 ? void 0 : _b.envVars) !== null && _d !== void 0 ? _d : {};
    if (Object.keys(envVars).length === 0) {
        return;
    }
    if ((_c = request.options) === null || _c === void 0 ? void 0 : _c.agentId) {
        throw new ConnectError("cloud.envVars cannot be used with a caller-supplied agentId", Code.InvalidArgument);
    }
    for (const name of Object.keys(envVars)) {
        if (name.startsWith("CURSOR_")) {
            throw new ConnectError("cloud.envVars names cannot start with CURSOR_", Code.InvalidArgument);
        }
    }
}
function parseOffset(offset) {
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
// proto3 cannot distinguish "apiKey field omitted" from `apiKey: ""`, so an
// empty key here means the SDK caller did not supply one. Refuse to fall back
// to the bridge's process env so misconfigured callers fail closed.
function requireCursorRequestApiKey(options) {
    const converted = protoCursorRequestOptionsToSdk(options);
    if (converted === undefined || !converted.apiKey) {
        throw new ConnectError("API key is required for cloud catalog calls.", Code.Unauthenticated);
    }
    return converted;
}
export function createSdkBridgeControlService(options = {}) {
    return {
        ping: async (_request) => new PingResponse({ message: "pong" }),
        shutdown: async (request) => {
            var _a;
            await ((_a = options.shutdown) === null || _a === void 0 ? void 0 : _a.call(options, request.graceSeconds));
            return new ShutdownResponse();
        },
        getVersion: async (_request) => new GetVersionResponse({
            bridgeVersion: CURSOR_SDK_BRIDGE_VERSION,
            protocolVersion: CURSOR_SDK_BRIDGE_PROTOCOL_VERSION,
            capabilities: CURSOR_SDK_BRIDGE_CAPABILITIES,
        }),
        // Point the bridge at a host custom-tool callback server after startup, so
        // clients attaching via Client.connect can register execute handlers. An
        // empty url clears any configured callback. Loopback/same-host only.
        setToolCallback: async (request) => {
            const url = request.url.trim();
            const authToken = request.authToken.trim();
            if (!url) {
                setBridgeToolCallbackConfig(undefined);
                resetBridgeToolCallbackClient();
                return new SetToolCallbackResponse();
            }
            if (!authToken) {
                throw new ConnectError("SetToolCallback requires an auth_token when url is set.", Code.InvalidArgument);
            }
            setBridgeToolCallbackConfig({ url, authToken });
            resetBridgeToolCallbackClient();
            return new SetToolCallbackResponse();
        },
    };
}
