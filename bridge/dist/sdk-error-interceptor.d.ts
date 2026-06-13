import { type Interceptor } from "@connectrpc/connect";
/**
 * Reshape a thrown error into a `ConnectError` with a recognized Connect code,
 * iff the error is a structured `CursorSdkError` with a code we know about.
 * Returns the input unchanged for any other error type so non-SDK errors
 * (raw `ConnectError`s, generic `Error`s) propagate as-is.
 */
export declare function maybeTranslateSdkError(error: unknown): unknown;
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
export declare function createSdkErrorInterceptor(): Interceptor;
//# sourceMappingURL=sdk-error-interceptor.d.ts.map