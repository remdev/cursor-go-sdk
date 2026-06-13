/**
 * Opt-in process-level "survive uncaughtException / unhandledRejection"
 * hooks for the bridge entrypoint.
 *
 * BACKGROUND
 *
 * The 2026-05-05 e2e suite at `bruce/sdk-python-e2e-suite-v2` showed the
 * cursor-sdk-bridge node subprocess dying mid-run (around the local-agent
 * test) and every subsequent cloud-touching test failing with
 * `connection refused`. The leak is somewhere in the handler chain —
 * suspects: local-agent disposal, KV teardown, run stream close — but the
 * exact path has not yet been identified.
 *
 * WHY THIS IS OPT-IN
 *
 * Per Node docs, the process is in an undefined state after
 * `uncaughtException`; surviving may leave the in-memory
 * `CursorSdkBridgeRegistry` or live IPC channels in a corrupted state.
 * Globally swallowing those signals is the textbook "On Error Resume Next"
 * anti-pattern: every future bridge bug gets silently absorbed instead of
 * surfacing in CI / dev / production.
 *
 * So this module **defaults to off**. Production, dev, and CI all see the
 * standard Node behavior — uncaughtException kills the process loudly.
 * The e2e suite (and only the e2e suite) opts in by setting
 * `CURSOR_SDK_BRIDGE_SURVIVE_UNCAUGHT=1`, with the understanding that the
 * survivor is a temporary band-aid until the leaking handler is fixed.
 *
 * REMOVAL CRITERIA
 *
 * Once the underlying handler is identified and hardened, this module
 * (and the env var) should be deleted. Tracked in the cursor-sdk roadmap.
 */
export const SURVIVE_UNCAUGHT_ENV_VAR = "CURSOR_SDK_BRIDGE_SURVIVE_UNCAUGHT";
/**
 * Returns true if the env signals that the caller has consciously opted
 * into the "On Error Resume Next" survivor handlers. Truthy values that
 * count as opt-in: `"1"`, `"true"`, `"yes"` (case-insensitive). Anything
 * else — including `undefined`, `""`, `"0"`, `"false"`, `"no"` — is
 * treated as "do not install survivors" so the process fails loud on
 * uncaughtException, as Node intends.
 */
export function shouldInstallProcessErrorSurvivors(env) {
    const raw = env[SURVIVE_UNCAUGHT_ENV_VAR];
    if (raw === undefined) {
        return false;
    }
    return /^(1|true|yes)$/i.test(raw.trim());
}
/**
 * Installs `uncaughtException` and `unhandledRejection` handlers that log
 * the offending stack to stderr and let the event loop continue. **Only
 * call this when `shouldInstallProcessErrorSurvivors` returns true.**
 *
 * Safe to call once per process. Calling twice will register the same
 * handler shape twice, which is harmless but wastes a process listener
 * slot (and may trip Node's `MaxListenersExceededWarning` if combined
 * with other registrations).
 */
export function installProcessErrorSurvivors() {
    process.on("uncaughtException", error => {
        process.stderr.write(`cursor-sdk-bridge uncaughtException (surviving): ${error instanceof Error ? error.stack || error.message : String(error)}\n`);
    });
    process.on("unhandledRejection", reason => {
        process.stderr.write(`cursor-sdk-bridge unhandledRejection (surviving): ${reason instanceof Error
            ? reason.stack || reason.message
            : String(reason)}\n`);
    });
}
