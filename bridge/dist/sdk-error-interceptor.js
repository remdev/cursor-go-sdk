import { SdkErrorCode, SdkErrorDetails, } from "@anysphere/proto/sdk/v1/sdk_errors_pb.js";
import { Code, ConnectError, } from "@connectrpc/connect";
import { CursorSdkError } from "@cursor/sdk";
// Map structured `CursorSdkError.code` strings to Connect codes. Without this
// translation, connect-node's adapter collapses every non-`ConnectError` into
// `Code.Internal`, hiding the structured code from SDK consumers.
const SDK_ERROR_CODE_TO_CONNECT_CODE = {
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
const SDK_ERROR_NAME_TO_CONNECT_CODE = {
    AuthenticationError: Code.Unauthenticated,
    ConfigurationError: Code.InvalidArgument,
    AgentBusyError: Code.FailedPrecondition,
    RateLimitError: Code.ResourceExhausted,
    AgentNotFoundError: Code.NotFound,
};
/**
 * Reshape a thrown error into a `ConnectError` with a recognized Connect code,
 * iff the error is a structured `CursorSdkError` with a code we know about.
 * Returns the input unchanged for any other error type so non-SDK errors
 * (raw `ConnectError`s, generic `Error`s) propagate as-is.
 */
export function maybeTranslateSdkError(error) {
    var _a;
    if (error instanceof CursorSdkError) {
        const connectCode = (_a = (error.code ? SDK_ERROR_CODE_TO_CONNECT_CODE[error.code] : undefined)) !== null && _a !== void 0 ? _a : SDK_ERROR_NAME_TO_CONNECT_CODE[error.name];
        if (connectCode !== undefined) {
            return new ConnectError(error.message, connectCode, undefined, [sdkErrorDetailsFromSdkError(error)], error);
        }
    }
    return error;
}
function sdkErrorDetailsFromSdkError(error) {
    const metadata = error;
    return new SdkErrorDetails({
        requestId: error.requestId,
        sdkErrorCode: sdkErrorCodeFromString(error.code),
        message: error.message,
        helpUrl: metadata.helpUrl,
        provider: metadata.provider,
    });
}
function sdkErrorCodeFromString(code) {
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
/**
 * Server-side Connect interceptor that funnels SDK errors thrown anywhere in
 * the bridge's RPC handlers through `maybeTranslateSdkError`. Replaces the
 * per-handler `translateSdkErrorToConnectError(...)` wrappers that used to
 * decorate every unary and streaming RPC.
 *
 * The Connect interceptor pipeline runs in both directions:
 * 1. For unary RPCs, errors thrown by the handler reject the awaited
 *    `next(req)` call directly; the catch translates and rethrows.
 * 2. For streaming RPCs, the response is a `StreamResponse` whose
 *    `.message` is an `AsyncIterable<O>`; errors thrown mid-stream surface
 *    when the consumer iterates, so we re-wrap the iterable to translate
 *    on every `for await` step.
 *
 * Errors that are already `ConnectError` (e.g. raw `Code.InvalidArgument`
 * argument validation thrown by handlers) pass through unchanged.
 */
export function createSdkErrorInterceptor() {
    return next => async req => {
        let response;
        try {
            response = await next(req);
        }
        catch (error) {
            throw maybeTranslateSdkError(error);
        }
        if (!response.stream) {
            return response;
        }
        return Object.assign(Object.assign({}, response), { message: translateStreamErrors(response.message) });
    };
}
async function* translateStreamErrors(inner) {
    try {
        yield* inner;
    }
    catch (error) {
        throw maybeTranslateSdkError(error);
    }
}
