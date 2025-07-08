async function verifyToken(r) {
    try {
        const secret = r.variables.jwt_secret;
        if (!secret) {
            r.error("JWT secret key is not set in Nginx variables.");
            r.return(500, "Internal Server Configuration Error");
            return;
        }

        const authHeader = r.headersIn['Authorization'];
        if (!authHeader) {
            r.return(401, "Authorization header is required");
            return;
        }

        const token = authHeader.split(' ')[1];
        if (!token) {
            r.return(401, "Bearer token is missing");
            return;
        }

        // Import thư viện jwt của NJS
        const jwt = require('jsonwebtoken');

        // Xác thực token
        const claims = await jwt.verify(token, secret);
        
        // Gửi các header đã được xác thực đến upstream service
        r.headersOut['X-User-ID'] = claims.UserId;
        r.headersOut['X-User-Role'] = claims.Role;
        r.headersOut['X-User-Email'] = claims.Email;
        
        // Trả về 200 để cho phép request đi tiếp
        r.return(200);

    } catch (e) {
        r.error(`JWT verification failed: ${e.message}`);
        r.return(401, `Unauthorized: ${e.message}`);
    }
}

export default { verifyToken };