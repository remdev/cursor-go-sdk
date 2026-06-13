import type { Run, SDKAgent } from "@cursor/sdk";
export declare class CursorSdkBridgeRegistry {
    private readonly agents;
    private readonly cloudAgentIds;
    private readonly runs;
    registerAgent(agent: SDKAgent, options?: {
        cloud?: boolean;
    }): Promise<void>;
    getAgent(agentId: string): SDKAgent | undefined;
    isCloudAgent(agentId: string): boolean;
    disposeAgent(agentId: string): Promise<void>;
    registerRun(run: Run): void;
    getRun(runId: string): Run | undefined;
    registerRuns(runs: Run[]): void;
    dispose(): Promise<void>;
    private agentIdsFor;
    private deleteAgentReferences;
}
//# sourceMappingURL=registry.d.ts.map