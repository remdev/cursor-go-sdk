import { AgentInfoStatus, AgentMessage, AgentModeOption, AgentDefinition, CloudAgentInfo, CloudEnvironment, CloudEnvironmentType, CustomToolDefinition, DownloadArtifactChunk, HttpMcpTransportType, LocalAgentInfo, ModelParameterDefinition, ModelParameterDefinitionValue, ModelParameterValue, ModelSelection, ModelVariant, ConversationStep as ProtoConversationStep, InteractionUpdate as ProtoInteractionUpdate, RunGitBranchInfo, RunGitInfo, RunLifecycleStatus, RunResult, RunSnapshot, RunStreamDone, RunStreamMessage, RunStreamResult, Runtime, SdkAgentInfo, SdkArtifact, SdkMessage, SdkModel, SdkRepository, SdkUser, SettingSource, } from "../gen/ts/sdk/v1/sdk_messages_pb.js";
import { Struct, Timestamp, type PlainMessage } from "@bufbuild/protobuf";
import { attachHostCustomToolsToSdkLocalOptions, protoCustomToolDefinitionsToHost, } from "./bridge-custom-tools.js";
import { createLocalAgentStoreFromProtoConfig } from "./bridge-local-agent-store.js";
import type { ListAgentsOptions } from "@cursor/sdk";
export function protoAgentOptionsToSdk(options, defaultCwd, customToolsAgentId) {
    if (!options) {
        // No options at all -> a bare local create/resume; default its cwd to the
        // launch workspace so it lands where reads (which also default to the
        // workspace) will find it.
        return defaultCwd ? { local: { cwd: defaultCwd } } : {};
    }
    return stripUndefined({
        model: protoModelSelectionToSdk(options.model),
        apiKey: emptyToUndefined(options.apiKey),
        name: emptyToUndefined(options.name),
        // Only default the cwd for local agents; cloud agents never carry a cwd.
        local: protoLocalAgentOptionsToSdk(options.local, options.cloud ? undefined : defaultCwd, customToolsAgentId),
        cloud: protoCloudAgentOptionsToSdk(options.cloud),
        mcpServers: protoMcpServersToSdk(options.mcpServers),
        agents: protoAgentDefinitionsToSdk(options.agents),
        agentId: emptyToUndefined(options.agentId),
        mode: protoAgentModeOptionToSdk(options.mode),
    });
}
export function protoSendOptionsToSdk(options) {
    if (!options) {
        return undefined;
    }
    const localBase = options.local
        ? stripUndefined({ force: options.local.force })
        : undefined;
    return stripUndefined({
        model: protoModelSelectionToSdk(options.model),
        mcpServers: protoMcpServersToSdk(options.mcpServers),
        local: localBase,
        mode: protoAgentModeOptionToSdk(options.mode),
    });
}
function protoAgentModeOptionToSdk(mode) {
    switch (mode) {
        case AgentModeOption.AGENT:
            return "agent";
        case AgentModeOption.PLAN:
            return "plan";
        default:
            return undefined;
    }
}
export function protoUserMessageToSdk(message) {
    return {
        text: message?.text ?? "",
        images: message?.images && message.images.length > 0
            ? message.images.map(protoImageToSdk)
            : undefined,
    };
}
export function sdkAgentToCreateResponse(agent) {
    return {
        agentId: agent.agentId,
        model: sdkModelSelectionToProto(agent.model),
    };
}
export function sdkRunToSnapshot(run) {
    return new RunSnapshot({
        runId: run.id,
        agentId: run.agentId,
        status: sdkRunStatusToProto(run.status),
        result: run.result ?? "",
        model: sdkModelSelectionToProto(run.model),
        durationMs: numberToProtoUint64(run.durationMs),
        git: sdkGitToProto(run.git),
        createdAt: timestampFromMillis(run.createdAt),
    });
}
export function sdkRunResultToProto(result, agentId = "", createdAtMs?) {
    return new RunResult({
        runId: result.id,
        agentId,
        status: sdkRunStatusToProto(result.status),
        result: result.result ?? "",
        model: sdkModelSelectionToProto(result.model),
        durationMs: numberToProtoUint64(result.durationMs),
        git: sdkGitToProto(result.git),
        createdAt: timestampFromMillis(createdAtMs),
    });
}
export function sdkMessageToRunStreamMessage(message, offset) {
    return new RunStreamMessage({
        envelope: {
            case: "sdkMessage",
            value: new SdkMessage({
                type: getMessageType(message),
                message: valueToStruct(message),
            }),
        },
        offset,
    });
}
export function sdkRunResultToRunStreamMessage(run, result, offset) {
    return new RunStreamMessage({
        envelope: {
            case: "result",
            value: new RunStreamResult({
                agentId: run.agentId,
                runId: run.id,
                status: sdkRunStatusToProto(result.status),
                result: sdkRunResultToProto(result, run.agentId, run.createdAt),
            }),
        },
        offset,
    });
}
export function sdkRunDoneToRunStreamMessage(run, offset) {
    return new RunStreamMessage({
        envelope: {
            case: "done",
            value: new RunStreamDone({ agentId: run.agentId, runId: run.id }),
        },
        offset,
    });
}
export function sdkInteractionUpdateToRunStreamMessage(update, offset) {
    return new RunStreamMessage({
        envelope: {
            case: "interactionUpdate",
            value: new ProtoInteractionUpdate({
                type: getPayloadType(update),
                update: valueToStruct(update),
            }),
        },
        offset,
    });
}
export function sdkConversationStepToRunStreamMessage(step, offset) {
    return new RunStreamMessage({
        envelope: {
            case: "step",
            value: new ProtoConversationStep({
                type: getPayloadType(step),
                step: valueToStruct(step),
            }),
        },
        offset,
    });
}
export function sdkListAgentsResponseItemsToProto(items) {
    return items.map(sdkAgentInfoToProto);
}
export function sdkAgentInfoToProto(info) {
    const runtimeInfo: PlainMessage<SdkAgentInfo>["runtimeInfo"] =
        info.runtime === "local"
            ? {
                case: "local",
                value: new LocalAgentInfo({ cwd: info.cwd ?? "" }),
            }
            : info.runtime === "cloud"
                ? {
                    case: "cloud",
                    value: new CloudAgentInfo({
                        env: sdkCloudEnvToProto(info.env),
                        repos: info.repos ?? [],
                    }),
                }
                : { case: undefined };
    return new SdkAgentInfo({
        agentId: info.agentId,
        name: info.name,
        summary: info.summary,
        lastModified: timestampFromMillis(info.lastModified),
        status: sdkAgentInfoStatusToProto(info.status),
        createdAt: timestampFromMillis(info.createdAt),
        archived: info.archived ?? false,
        runtimeInfo,
    });
}
export function sdkAgentMessagesToProto(messages) {
    return messages.map(message => new AgentMessage({
        type: message.type,
        uuid: message.uuid,
        agentId: message.agent_id,
        message: valueToStruct(message.message),
    }));
}
export function sdkArtifactsToProto(artifacts) {
    return artifacts.map(artifact => new SdkArtifact({
        path: artifact.path,
        sizeBytes: numberToProtoUint64(artifact.sizeBytes),
        updatedAt: artifact.updatedAt,
    }));
}
export function sdkArtifactBufferToChunks(buffer, chunkBytes) {
    const chunks = [];
    for (let offset = 0; offset < buffer.length; offset += chunkBytes) {
        chunks.push(new DownloadArtifactChunk({
            data: new Uint8Array(buffer.subarray(offset, offset + chunkBytes)),
        }));
    }
    return chunks;
}
export function sdkUserToProto(user) {
    return new SdkUser({
        apiKeyName: user.apiKeyName,
        userId: numberToProtoUint64(user.userId),
        userEmail: user.userEmail ?? "",
        userFirstName: user.userFirstName ?? "",
        userLastName: user.userLastName ?? "",
        createdAt: user.createdAt,
    });
}
export function sdkModelsToProto(items) {
    return items.map(sdkModelToProto);
}
export function sdkRepositoriesToProto(items) {
    return items.map(item => new SdkRepository({ url: item.url }));
}
export function protoListAgentsOptionsToSdk(options, defaultCwd): ListAgentsOptions | undefined {
    if (!options) {
        return defaultCwd ? { runtime: "local", cwd: defaultCwd } : undefined;
    }
    const common = stripUndefined({
        limit: zeroToUndefined(options.limit),
        cursor: emptyToUndefined(options.cursor),
        cwd:
            options.runtime === Runtime.CLOUD
                ? undefined
                : (emptyToUndefined(options.cwd) ?? defaultCwd),
    });
    switch (options.runtime) {
        case Runtime.CLOUD:
            return stripUndefined({
                ...common,
                runtime: "cloud",
                prUrl: emptyToUndefined(options.prUrl),
                includeArchived: options.includeArchived,
                apiKey: emptyToUndefined(options.apiKey),
            }) as ListAgentsOptions;
        case Runtime.LOCAL:
            return stripUndefined({ ...common, runtime: "local" }) as ListAgentsOptions;
        default:
            return defaultCwd
                ? (stripUndefined({ ...common, runtime: "local", cwd: defaultCwd }) as ListAgentsOptions)
                : (common as ListAgentsOptions | undefined);
    }
}
export function protoListRunsOptionsToSdk(options, defaultCwd) {
    if (!options) {
        return defaultCwd ? stripUndefined({ cwd: defaultCwd }) : undefined;
    }
    const common = stripUndefined({
        limit: zeroToUndefined(options.limit),
        cursor: emptyToUndefined(options.cursor),
        cwd:
            options.runtime === Runtime.CLOUD
                ? undefined
                : (emptyToUndefined(options.cwd) ?? defaultCwd),
    });
    switch (options.runtime) {
        case Runtime.CLOUD:
            return stripUndefined({
                ...common,
                runtime: "cloud",
                apiKey: emptyToUndefined(options.apiKey),
            });
        case Runtime.LOCAL:
            return stripUndefined({ ...common, runtime: "local" });
        default:
            return common;
    }
}
export function protoGetRunOptionsToSdk(options, defaultCwd) {
    if (!options) {
        return defaultCwd ? stripUndefined({ cwd: defaultCwd }) : undefined;
    }
    if (options.runtime === Runtime.CLOUD) {
        const agentId = emptyToUndefined(options.agentId);
        if (!agentId) {
            return undefined;
        }
        return stripUndefined({
            runtime: "cloud",
            agentId,
            apiKey: emptyToUndefined(options.apiKey),
        });
    }
    return stripUndefined({
        runtime: options.runtime === Runtime.LOCAL ? "local" : undefined,
        cwd: emptyToUndefined(options.cwd) ?? defaultCwd,
    });
}
export function protoAgentOperationOptionsToSdk(options, defaultCwd) {
    if (!options) {
        return defaultCwd ? stripUndefined({ cwd: defaultCwd }) : undefined;
    }
    return stripUndefined({
        cwd: emptyToUndefined(options.cwd) ?? defaultCwd,
        apiKey: emptyToUndefined(options.apiKey),
    });
}
export function protoAgentMessagesOptionsToSdk(options, defaultCwd) {
    if (!options) {
        return defaultCwd ? stripUndefined({ cwd: defaultCwd }) : undefined;
    }
    const common = stripUndefined({
        limit: zeroToUndefined(options.limit),
        offset: zeroToUndefined(options.offset),
    });
    switch (options.runtime) {
        case Runtime.CLOUD:
            return stripUndefined({
                ...common,
                runtime: "cloud",
                apiKey: emptyToUndefined(options.apiKey),
            });
        case Runtime.LOCAL:
            return stripUndefined({
                ...common,
                runtime: "local",
                cwd: emptyToUndefined(options.cwd) ?? defaultCwd,
            });
        default:
            return stripUndefined({
                ...common,
                cwd: emptyToUndefined(options.cwd) ?? defaultCwd,
            });
    }
}
export function protoCursorRequestOptionsToSdk(options) {
    if (!options) {
        return undefined;
    }
    return stripUndefined({ apiKey: emptyToUndefined(options.apiKey) });
}
function protoLocalAgentOptionsToSdk(local, defaultCwd, customToolsAgentId) {
    if (!local) {
        return defaultCwd ? { cwd: defaultCwd } : undefined;
    }
    const base = stripUndefined({
        cwd:
            local.cwd.length === 0
                ? defaultCwd
                : local.cwd.length === 1
                    ? local.cwd[0]
                    : local.cwd,
        settingSources:
            local.settingSources.length === 0
                ? undefined
                : local.settingSources
                    .map(protoSettingSourceToSdk)
                    .filter((source) => source !== undefined),
        sandboxOptions:
            local.sandboxOptions?.enabled === undefined
                ? undefined
                : { enabled: local.sandboxOptions.enabled },
        store: protoLocalAgentStoreToSdk(local.store),
        autoReview: local.autoReview === undefined ? undefined : local.autoReview,
    });
    return attachHostCustomToolsToSdkLocalOptions(base, protoCustomToolDefinitionsToHost(local.customTools), customToolsAgentId);
}
function protoLocalAgentStoreToSdk(store) {
    if (!store) {
        return undefined;
    }
    return createLocalAgentStoreFromProtoConfig(store);
}
function protoCloudAgentOptionsToSdk(cloud) {
    if (!cloud) {
        return undefined;
    }
    return stripUndefined({
        env: protoCloudEnvToSdk(cloud.env),
        repos: cloud.repos.map(repo => stripUndefined({
            url: repo.url,
            startingRef: emptyToUndefined(repo.startingRef),
            prUrl: emptyToUndefined(repo.prUrl),
        })),
        workOnCurrentBranch: cloud.workOnCurrentBranch,
        autoCreatePR: cloud.autoCreatePr,
        skipReviewerRequest: cloud.skipReviewerRequest,
        envVars: Object.keys(cloud.envVars).length === 0 ? undefined : cloud.envVars,
    });
}
function protoMcpServersToSdk(servers) {
    const entries = Object.entries(servers)
        .map(([name, config]) => [name, protoMcpServerToSdk(config)])
        .filter((entry) => entry[1] !== undefined);
    return entries.length === 0 ? undefined : Object.fromEntries(entries);
}
function protoMcpServerToSdk(config) {
    switch (config.config.case) {
        case "stdio":
            return stripUndefined({
                type: "stdio",
                command: config.config.value.command,
                args: config.config.value.args.length > 0
                    ? config.config.value.args
                    : undefined,
                env: Object.keys(config.config.value.env).length > 0
                    ? config.config.value.env
                    : undefined,
                cwd: emptyToUndefined(config.config.value.cwd),
            });
        case "http":
            return stripUndefined({
                type: config.config.value.type === HttpMcpTransportType.SSE
                    ? "sse"
                    : "http",
                url: config.config.value.url,
                headers: Object.keys(config.config.value.headers).length > 0
                    ? config.config.value.headers
                    : undefined,
                auth: config.config.value.auth
                    ? stripUndefined({
                        CLIENT_ID: config.config.value.auth.clientId,
                        CLIENT_SECRET: emptyToUndefined(config.config.value.auth.clientSecret),
                        scopes: config.config.value.auth.scopes.length > 0
                            ? config.config.value.auth.scopes
                            : undefined,
                    })
                    : undefined,
            });
        default:
            return undefined;
    }
}
function protoAgentDefinitionsToSdk(agents: { [key: string]: AgentDefinition } | undefined) {
    if (!agents) {
        return undefined;
    }
    const entries = Object.entries(agents).map(([name, agent]) => {
        const mcpServers = [];
        let inlineConfigIndex = 0;
        for (const server of agent.mcpServers) {
            if (server.value.case === "name") {
                mcpServers.push(server.value.value);
                continue;
            }
            if (server.value.case === "inlineConfig") {
                const config = protoMcpServerToSdk(server.value.value);
                if (config) {
                    inlineConfigIndex += 1;
                    mcpServers.push({ [`${name}-inline-${inlineConfigIndex}`]: config });
                }
            }
        }
        return [
            name,
            stripUndefined({
                description: agent.description,
                prompt: agent.prompt,
                model: agent.inheritModel
                    ? "inherit"
                    : protoModelSelectionToSdk(agent.model),
                mcpServers: mcpServers.length === 0 ? undefined : mcpServers,
            }),
        ];
    });
    return entries.length === 0 ? undefined : Object.fromEntries(entries);
}
function protoModelSelectionToSdk(model) {
    if (!model || model.id.length === 0) {
        return undefined;
    }
    return {
        id: model.id,
        params: model.params.length === 0
            ? undefined
            : model.params.map(param => ({ id: param.id, value: param.value })),
    };
}
function protoImageToSdk(image) {
    const dimension = image.dimension
        ? { width: image.dimension.width, height: image.dimension.height }
        : undefined;
    if (image.source.case === "data") {
        return stripUndefined({
            data: image.source.value.data,
            mimeType: image.source.value.mimeType,
            dimension,
        });
    }
    return stripUndefined({
        url: image.source.case === "url" ? image.source.value.url : "",
        dimension,
    });
}
function sdkModelSelectionToProto(model) {
    if (!model) {
        return undefined;
    }
    return new ModelSelection({
        id: model.id,
        params: (model.params ?? []).map((param) =>
            new ModelParameterValue({ id: param.id, value: param.value }),
        ),
    });
}
function sdkModelToProto(model) {
    return new SdkModel({
        id: model.id,
        displayName: model.displayName,
        description: model.description ?? "",
        parameters: (model.parameters ?? []).map((parameter) =>
            new ModelParameterDefinition({
                id: parameter.id,
                displayName: parameter.displayName ?? "",
                values: parameter.values.map((value) =>
                    new ModelParameterDefinitionValue({
                        value: value.value,
                        displayName: value.displayName ?? "",
                    }),
                ),
            }),
        ),
        variants: (model.variants ?? []).map((variant) =>
            new ModelVariant({
                params: variant.params.map((param) =>
                    new ModelParameterValue({ id: param.id, value: param.value }),
                ),
                displayName: variant.displayName,
                description: variant.description ?? "",
                isDefault: variant.isDefault ?? false,
            }),
        ),
    });
}
function sdkGitToProto(git) {
    if (!git) {
        return undefined;
    }
    return new RunGitInfo({
        branches: git.branches.map((branch) =>
            new RunGitBranchInfo({
                repoUrl: branch.repoUrl,
                branch: branch.branch ?? "",
                prUrl: branch.prUrl ?? "",
            }),
        ),
    });
}
function protoCloudEnvToSdk(env) {
    if (!env || env.type === CloudEnvironmentType.UNSPECIFIED) {
        return undefined;
    }
    const type = env.type === CloudEnvironmentType.POOL
        ? "pool"
        : env.type === CloudEnvironmentType.MACHINE
            ? "machine"
            : "cloud";
    return stripUndefined({ type, name: emptyToUndefined(env.name) });
}
function sdkCloudEnvToProto(env) {
    if (!env) {
        return undefined;
    }
    return new CloudEnvironment({
        type:
            env.type === "pool"
                ? CloudEnvironmentType.POOL
                : env.type === "machine"
                    ? CloudEnvironmentType.MACHINE
                    : CloudEnvironmentType.CLOUD,
        name: env.name ?? "",
    });
}
function protoSettingSourceToSdk(source) {
    switch (source) {
        case SettingSource.PROJECT:
            return "project";
        case SettingSource.USER:
            return "user";
        case SettingSource.TEAM:
            return "team";
        case SettingSource.MDM:
            return "mdm";
        case SettingSource.PLUGINS:
            return "plugins";
        case SettingSource.ALL:
            return "all";
        default:
            return undefined;
    }
}
function sdkRunStatusToProto(status) {
    switch (status) {
        case "running":
            return RunLifecycleStatus.RUNNING;
        case "finished":
            return RunLifecycleStatus.FINISHED;
        case "error":
            return RunLifecycleStatus.ERROR;
        case "cancelled":
            return RunLifecycleStatus.CANCELLED;
        default:
            return RunLifecycleStatus.UNSPECIFIED;
    }
}
function sdkAgentInfoStatusToProto(status) {
    switch (status) {
        case "running":
            return AgentInfoStatus.RUNNING;
        case "finished":
            return AgentInfoStatus.FINISHED;
        case "error":
            return AgentInfoStatus.ERROR;
        default:
            return AgentInfoStatus.UNSPECIFIED;
    }
}
function getMessageType(message) {
    return getPayloadType(message);
}
function getPayloadType(message) {
    return typeof message === "object" &&
        message !== null &&
        "type" in message &&
        typeof message.type === "string"
        ? message.type
        : "unknown";
}
function valueToStruct(value) {
    try {
        return Struct.fromJson(toJsonValue(value));
    }
    catch {
        return Struct.fromJson({});
    }
}
function toJsonValue(value) {
    switch (typeof value) {
        case "string":
        case "number":
        case "boolean":
            return value;
        case "bigint":
            return value.toString();
        case "object":
            if (value === null) {
                return null;
            }
            if (Array.isArray(value)) {
                return value.map(item => item === undefined ? null : toJsonValue(item));
            }
            return Object.fromEntries(Object.entries(value)
                .filter(([, nested]) => nested !== undefined)
                .map(([key, nested]) => [key, toJsonValue(nested)]));
        default:
            throw new TypeError(`Cannot encode ${typeof value} in protobuf Struct`);
    }
}
function timestampFromMillis(value) {
    if (value === undefined || !Number.isFinite(value) || value <= 0) {
        return undefined;
    }
    return Timestamp.fromDate(new Date(value));
}
function numberToProtoUint64(value) {
    if (value === undefined || !Number.isFinite(value)) {
        return BigInt(0);
    }
    return BigInt(Math.max(0, Math.trunc(value)));
}
function zeroToUndefined(value) {
    return value === 0 ? undefined : value;
}
function emptyToUndefined(value) {
    return value.length === 0 ? undefined : value;
}
function stripUndefined(value) {
    return Object.fromEntries(Object.entries(value).filter(([, entry]) => entry !== undefined));
}
