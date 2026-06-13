import { randomBytes, timingSafeEqual } from "node:crypto";
import { Code, ConnectError } from "@connectrpc/connect";
export function createBridgeAuthToken() {
    return randomBytes(32).toString("base64url");
}
export function createBridgeAuthInterceptor(authToken) {
    return next => async request => {
        const authorization = request.header.get("authorization");
        const bearer = (authorization === null || authorization === void 0 ? void 0 : authorization.startsWith("Bearer "))
            ? authorization.slice("Bearer ".length)
            : undefined;
        if (!bearer || !constantTimeEquals(bearer, authToken)) {
            throw new ConnectError("Unauthorized", Code.Unauthenticated);
        }
        return await next(request);
    };
}
function constantTimeEquals(left, right) {
    const leftBuffer = Buffer.from(left);
    const rightBuffer = Buffer.from(right);
    return (leftBuffer.length === rightBuffer.length &&
        timingSafeEqual(leftBuffer, rightBuffer));
}
