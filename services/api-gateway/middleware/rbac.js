const { roles, permissions, publicRoutes } = require('../rbac_config');

// --- Optimization: Pre-compile permission routes ---
const permissionMatchers = {};
for (const role in permissions) {
    permissionMatchers[role] = permissions[role].map(p => {
        if (p === '*') return '*';
        return {
            ...p,
            regex: new RegExp(`^${p.path.replace(/:\w+/g, '[^/]+')}$`),
        };
    });
}
// --- End Optimization ---

/**
 * Middleware for Role-Based Access Control (RBAC).
 * It checks if the user's role has the necessary permission to access a route.
 * It also handles ownership checks for resources.
 */
const rbacMiddleware = (req, res, next) => {
    // The auth middleware has already identified public routes.
    if (req.isPublic) {
        return next();
    }

    // If the route is not public, a user must be attached to the request.
    // The auth middleware should have already sent a 401 if not.
    if (!req.user) {
        return res.status(401).json({ error: 'Authentication required.' });
    }

    const role = req.user.role;
    const userPermissions = permissionMatchers[role] || [];

    // Admins have universal access.
    if (userPermissions.includes('*')) {
        return next();
    }

    // Find a permission that matches the requested route and method.
    const matchedPermission = userPermissions.find(p =>
        p.regex && p.regex.test(req.path) && p.method === req.method
    );

    if (matchedPermission) {
        // If the permission requires ownership, perform a check.
        if (matchedPermission.own === true) {
            const requesterId = req.user.user_id.toString();
            const resourceIdParam = matchedPermission.param;

            // If the permission defines a URL parameter for the resource ID, check it.
            if (resourceIdParam && req.params[resourceIdParam]) {
                if (requesterId === req.params[resourceIdParam]) {
                    return next(); // User is accessing their own resource.
                }
                // If the IDs don't match, access is forbidden.
                return res.status(403).json({ error: 'Forbidden: You can only access your own resources.' });
            }
            // If no param is specified (e.g., /api/users/profile), it's a direct resource of the user.
            return next();
        }
        // If `own` is not true, the permission is granted without an ownership check.
        return next();
    }

    // If no matching permission was found, deny access.
    return res.status(403).json({ error: 'Forbidden: You do not have permission to access this resource.' });
};

module.exports = rbacMiddleware;
