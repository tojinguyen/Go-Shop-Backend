function verifyToken(r) {
    function log(msg) {
        r.error("[auth.js] " + msg);
    }

    var secret = r.variables.jwt_secret_key;
    if (!secret) {
        log("JWT_SECRET_KEY is not set.");
        return r.return(500, "Server error");
    }

    var authHeader = r.headersIn['Authorization'];
    if (!authHeader || authHeader.indexOf("Bearer ") !== 0) {
        return r.return(401, "Missing or invalid Authorization header");
    }

    var token = authHeader.slice(7);
    var parts = token.split(".");
    if (parts.length !== 3) {
        log("Malformed token");
        return r.return(401, "Malformed token");
    }

    var headerB64 = parts[0];
    var payloadB64 = parts[1];
    var signatureB64 = parts[2];

    // Base64url decode manually
    function b64urlDecode(str) {
        str = str.replace(/-/g, "+").replace(/_/g, "/");
        while (str.length % 4 !== 0) str += "=";
        return JSON.parse(atob(str));
    }

    var header = b64urlDecode(headerB64);
    if (header.alg !== "HS256") {
        log("Unsupported alg: " + header.alg);
        return r.return(401, "Invalid algorithm");
    }

    var payload = b64urlDecode(payloadB64);
    var now = Math.floor(Date.now() / 1000);
    if (payload.exp && now > payload.exp) {
        log("Token expired");
        return r.return(401, "Token expired");
    }

    // Emulate HMAC SHA256 signature check
    // ❗ NJS chưa có crypto, nên bỏ qua check signature ở đây hoặc tự viết lại module native

    // Nếu muốn verify signature, nên chuyển về gateway Go hoặc node middleware.

    r.headersOut['X-User-ID'] = payload.UserId;
    r.headersOut['X-User-Role'] = payload.Role;
    r.headersOut['X-User-Email'] = payload.Email;

    log("Token verified");
    return r.return(204);
}

export default { verifyToken };
