import type { CursorSdkBridgeRegistry } from "./registry.js";

export interface CliArgs {
    help?: boolean;
    host?: string;
    port?: string;
    workspace?: string;
    stateRoot?: string;
    toolCallbackUrl?: string;
    toolCallbackAuthToken?: string;
    localStore?: string;
    maxConcurrentAgents?: number;
    maxMessageBytes?: number;
}

export interface BridgeServerAddress {
    schemaVersion: number;
    serverVersion: string;
    pid: number;
    transport: string;
    protocol: string;
    host: string;
    port: number;
    url: string;
    authToken: string;
    workspaceRef: string;
    stateRoot: string;
    maxConcurrentAgents?: number;
    maxMessageBytes?: number;
}

export interface StartCursorSdkBridgeServerOptions {
    host?: string;
    port?: number;
    allowNonLoopbackHost?: boolean;
    authToken?: string;
    registry?: CursorSdkBridgeRegistry;
    workspaceRef?: string;
    stateRoot?: string;
    maxConcurrentAgents?: number;
    maxMessageBytes?: number;
    writeReady?: (address: BridgeServerAddress) => void;
}

export interface CursorSdkBridgeServiceOptions {
    registry?: CursorSdkBridgeRegistry;
    workspaceRef?: string;
}

export interface SdkBridgeControlServiceOptions {
    shutdown?: (graceSeconds: number) => void | Promise<void>;
}

export interface MountSdkBridgeRoutesOptions {
    registry: CursorSdkBridgeRegistry;
    workspaceRef: string;
    shutdown: (graceSeconds: number) => void | Promise<void>;
}
