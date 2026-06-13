export class CursorSdkBridgeRegistry {
    agents = new Map();
    cloudAgentIds = new Set();
    runs = new Map();
    async registerAgent(agent, options) {
        const existingAgent = this.agents.get(agent.agentId);
        try {
            if (existingAgent && existingAgent !== agent) {
                await existingAgent[Symbol.asyncDispose]();
            }
        }
        finally {
            if (existingAgent && existingAgent !== agent) {
                this.deleteAgentReferences(existingAgent);
            }
            this.agents.set(agent.agentId, agent);
            if (options === null || options === void 0 ? void 0 : options.cloud) {
                this.cloudAgentIds.add(agent.agentId);
            }
            else {
                this.cloudAgentIds.delete(agent.agentId);
            }
        }
    }
    getAgent(agentId) {
        return this.agents.get(agentId);
    }
    isCloudAgent(agentId) {
        return this.cloudAgentIds.has(agentId);
    }
    async disposeAgent(agentId) {
        const agent = this.agents.get(agentId);
        const removedAgentIds = agent
            ? this.agentIdsFor(agent)
            : new Set([agentId]);
        try {
            if (agent) {
                await agent[Symbol.asyncDispose]();
            }
        }
        finally {
            for (const registeredAgentId of removedAgentIds) {
                this.agents.delete(registeredAgentId);
                this.cloudAgentIds.delete(registeredAgentId);
            }
            for (const [runId, run] of this.runs.entries()) {
                if (removedAgentIds.has(run.agentId)) {
                    this.runs.delete(runId);
                }
            }
        }
    }
    registerRun(run) {
        this.runs.set(run.id, run);
    }
    getRun(runId) {
        return this.runs.get(runId);
    }
    registerRuns(runs) {
        for (const run of runs) {
            this.registerRun(run);
        }
    }
    async dispose() {
        const errors = [];
        try {
            for (const agent of new Set(this.agents.values())) {
                try {
                    await agent[Symbol.asyncDispose]();
                }
                catch (err) {
                    errors.push(err);
                }
            }
        }
        finally {
            this.agents.clear();
            this.cloudAgentIds.clear();
            this.runs.clear();
        }
        if (errors.length === 1) {
            throw errors[0];
        }
        if (errors.length > 1) {
            throw new AggregateError(errors, "Failed to dispose bridge agents");
        }
    }
    agentIdsFor(agent) {
        const agentIds = new Set();
        for (const [agentId, registeredAgent] of this.agents.entries()) {
            if (registeredAgent === agent) {
                agentIds.add(agentId);
            }
        }
        return agentIds;
    }
    deleteAgentReferences(agent) {
        for (const agentId of this.agentIdsFor(agent)) {
            this.agents.delete(agentId);
            this.cloudAgentIds.delete(agentId);
        }
    }
}
