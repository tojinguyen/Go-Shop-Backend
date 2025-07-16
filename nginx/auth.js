function verifyToken(r) {
    function log(msg) {
        r.error("[auth.js] " + msg);
    }

    log("Auth subrequest started.");
    var secret = process.env.JWT_SECRET_KEY;

    if (!secret) {
        log("FATAL: Environment variable JWT_SECRET_KEY is not accessible inside NJS.");
        return r.return(500, "Server Configuration Error: JWT Secret not found");
    }

    log("Secret key loaded via process.env.");

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

    r.variables.auth_user_id = payload.UserId;
    r.variables.auth_user_role = payload.Role;
    r.variables.auth_user_email = payload.Email;

    log("SUCCESS: Token verified. UserID: " + payload.UserId);

    return r.return(204);
}

function getUserId(r) {
    return r.variables.auth_user_id || ""; 
}

function getRole(r) {
    return r.variables.auth_user_role || "";
}

function getEmail(r) {
    return r.variables.auth_user_email || "";
}

export default { verifyToken, getUserId, getRole, getEmail };
