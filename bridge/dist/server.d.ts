import { CursorSdkBridgeRegistry } from "./registry.js";
export interface StartBridgeServerOptions {
    host?: string;
    port?: number;
    registry?: CursorSdkBridgeRegistry;
    authToken?: string;
    allowNonLoopbackHost?: boolean;
    workspaceRef?: string;
    stateRoot?: string;
    maxConcurrentAgents?: number;
    maxMessageBytes?: number;
    writeReady?: (address: BridgeServerAddress) => void;
}
export interface BridgeServerAddress {
    schemaVersion: 1;
    serverVersion: string;
    pid: number;
    transport: "tcp";
    protocol: "connect";
    host: string;
    port: number;
    url: string;
    authToken: string;
    authTokenFile?: string;
    workspaceRef: string;
    stateRoot: string;
    maxConcurrentAgents?: number;
    maxMessageBytes?: number;
}
export interface BridgeServerHandle {
    address: BridgeServerAddress;
    close(): Promise<void>;
}
export declare function startCursorSdkBridgeServer(options?: StartBridgeServerOptions): Promise<BridgeServerHandle>;
export declare const startBridgeServer: typeof startCursorSdkBridgeServer;
//# sourceMappingURL=server.d.ts.map