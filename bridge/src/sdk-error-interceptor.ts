import { SdkErrorCode, SdkErrorDetails } from "../gen/ts/sdk/v1/sdk_errors_pb.js";
import { Code, ConnectError, type Interceptor } from "@connectrpc/connect";
import { CursorSdkError } from "@cursor/sdk";

interface SdkErrorMetadata {
    helpUrl?: string;
    provider?: string;
}

const SDK_ERROR_CODE_TO_CONNECT_CODE: Record<string, Code> = {
    invalid_model: Code.InvalidArgument,
    unauthenticated: Code.Unauthenticated,
    permission_denied: Code.PermissionDenied,
    not_found: Code.NotFound,
    agent_not_found: Code.NotFound,
    run_not_found: Code.NotFound,
    unknown_agent: Code.NotFound,
    unknown_run: Code.NotFound,
    rate_limited: Code.ResourceExhausted,
    stream_buffer_overflow: Code.ResourceExhausted,
    integration_not_connected: Code.FailedPrecondition,
    repository_access: Code.FailedPrecondition,
    agent_busy: Code.FailedPrecondition,
};

const SDK_ERROR_NAME_TO_CONNECT_CODE: Record<string, Code> = {
    AuthenticationError: Code.Unauthenticated,
    ConfigurationError: Code.InvalidArgument,
    AgentBusyError: Code.FailedPrecondition,
    RateLimitError: Code.ResourceExhausted,
    AgentNotFoundError: Code.NotFound,
};

export function maybeTranslateSdkError(error: unknown): unknown {
    if (error instanceof CursorSdkError) {
        const connectCode =
            (error.code ? SDK_ERROR_CODE_TO_CONNECT_CODE[error.code] : undefined) ??
            SDK_ERROR_NAME_TO_CONNECT_CODE[error.name];
        if (connectCode !== undefined) {
            return new ConnectError(
                error.message,
                connectCode,
                undefined,
                [sdkErrorDetailsFromSdkError(error)],
                error,
            );
        }
    }
    return error;
}

function sdkErrorDetailsFromSdkError(error: CursorSdkError): SdkErrorDetails {
    const metadata = error as CursorSdkError & SdkErrorMetadata;
    return new SdkErrorDetails({
        requestId: error.requestId,
        sdkErrorCode: sdkErrorCodeFromString(error.code),
        message: error.message,
        helpUrl: metadata.helpUrl,
        provider: metadata.provider,
    });
}

function sdkErrorCodeFromString(code: string | undefined): SdkErrorCode {
    switch (code) {
        case "unauthenticated":
            return SdkErrorCode.UNAUTHORIZED;
        case "invalid_model":
            return SdkErrorCode.INVALID_MODEL;
        case "permission_denied":
            return SdkErrorCode.ROLE_FORBIDDEN;
        case "agent_not_found":
        case "unknown_agent":
            return SdkErrorCode.AGENT_NOT_FOUND;
        case "run_not_found":
        case "unknown_run":
            return SdkErrorCode.RUN_NOT_FOUND;
        case "rate_limited":
            return SdkErrorCode.RATE_LIMIT_EXCEEDED;
        case "integration_not_connected":
        case "repository_access":
            return SdkErrorCode.REPOSITORY_ACCESS;
        case "agent_busy":
            return SdkErrorCode.AGENT_BUSY;
        default:
            return SdkErrorCode.UNSPECIFIED;
    }
}

export function createSdkErrorInterceptor(): Interceptor {
    return (next) => async (req) => {
        try {
            const response = await next(req);
            if (!response.stream) {
                return response;
            }
            return {
                ...response,
                message: translateStreamErrors(response.message),
            };
        } catch (error) {
            throw maybeTranslateSdkError(error);
        }
    };
}

async function* translateStreamErrors<T>(inner: AsyncIterable<T>): AsyncGenerator<T> {
    try {
        yield* inner;
    } catch (error) {
        throw maybeTranslateSdkError(error);
    }
}
