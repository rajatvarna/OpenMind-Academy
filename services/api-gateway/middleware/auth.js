const jwt = require('jsonwebtoken');
const fs = require('fs');
const { publicRoutes } = require('../rbac_config');

// --- Optimization: Pre-compile routes and read key on startup ---
const publicKey = fs.readFileSync(process.env.JWT_PUBLIC_KEY_PATH || '../secrets/jwtRS256.key.pub');
const publicRouteMatchers = publicRoutes.map(route => ({
    method: route.method,
    regex: new RegExp(`^${route.path.replace(/:\w+/g, '[^/]+')}$`),
}));
// --- End Optimization ---

/**
 * Middleware to handle JWT-based authentication.
 * It verifies the token from the Authorization header and attaches the decoded payload to req.user.
 * It also identifies public routes and skips authentication for them.
 */
const authMiddleware = (req, res, next) => {
    // Check if the request path matches any of the pre-compiled public routes.
    const isPublic = publicRouteMatchers.some(route =>
        route.method === req.method && route.regex.test(req.path)
    );

    // If the route is public, skip authentication and set a flag for the RBAC middleware.
    if (isPublic) {
        req.isPublic = true;
        return next();
    }

    // For protected routes, expect a 'Bearer' token in the Authorization header.
    const authHeader = req.headers.authorization;
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return res.status(401).json({ error: 'Authorization header is required' });
    }

    const token = authHeader.split(' ')[1];
    try {
        // Verify the token using the pre-loaded public key.
        const decoded = jwt.verify(token, publicKey, { algorithms: ['RS256'] });
        req.user = decoded; // Attach decoded user info (id, role) to the request object.
        next();
    } catch (err) {
        // Handle errors like expired tokens or invalid signatures.
        return res.status(401).json({ error: 'Invalid or expired token' });
    }
};

module.exports = authMiddleware;
