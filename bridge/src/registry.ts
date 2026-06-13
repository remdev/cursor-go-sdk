import type { Run, SDKAgent } from "@cursor/sdk";

export interface RegisterAgentOptions {
    cloud?: boolean;
}

export class CursorSdkBridgeRegistry {
    agents = new Map<string, SDKAgent>();
    cloudAgentIds = new Set<string>();
    runs = new Map<string, Run>();

    async registerAgent(agent: SDKAgent, options?: RegisterAgentOptions): Promise<void> {
        const existingAgent = this.agents.get(agent.agentId);
        try {
            if (existingAgent && existingAgent !== agent) {
                await existingAgent[Symbol.asyncDispose]();
            }
        } finally {
            if (existingAgent && existingAgent !== agent) {
                this.deleteAgentReferences(existingAgent);
            }
            this.agents.set(agent.agentId, agent);
            if (options?.cloud) {
                this.cloudAgentIds.add(agent.agentId);
            } else {
                this.cloudAgentIds.delete(agent.agentId);
            }
        }
    }

    getAgent(agentId: string): SDKAgent | undefined {
        return this.agents.get(agentId);
    }

    isCloudAgent(agentId: string): boolean {
        return this.cloudAgentIds.has(agentId);
    }

    async disposeAgent(agentId: string): Promise<void> {
        const agent = this.agents.get(agentId);
        const removedAgentIds = agent ? this.agentIdsFor(agent) : new Set([agentId]);
        try {
            if (agent) {
                await agent[Symbol.asyncDispose]();
            }
        } finally {
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

    registerRun(run: Run): void {
        this.runs.set(run.id, run);
    }

    getRun(runId: string): Run | undefined {
        return this.runs.get(runId);
    }

    registerRuns(runs: Run[]): void {
        for (const run of runs) {
            this.registerRun(run);
        }
    }

    async dispose(): Promise<void> {
        const errors: unknown[] = [];
        try {
            for (const agent of new Set(this.agents.values())) {
                try {
                    await agent[Symbol.asyncDispose]();
                } catch (err) {
                    errors.push(err);
                }
            }
        } finally {
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

    private agentIdsFor(agent: SDKAgent): Set<string> {
        const agentIds = new Set<string>();
        for (const [agentId, registeredAgent] of this.agents.entries()) {
            if (registeredAgent === agent) {
                agentIds.add(agentId);
            }
        }
        return agentIds;
    }

    private deleteAgentReferences(agent: SDKAgent): void {
        for (const agentId of this.agentIdsFor(agent)) {
            this.agents.delete(agentId);
            this.cloudAgentIds.delete(agentId);
        }
    }
}
