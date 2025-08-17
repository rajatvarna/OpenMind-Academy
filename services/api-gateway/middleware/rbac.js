const { roles, permissions, publicRoutes } = require('../rbac_config');

const rbacMiddleware = (req, res, next) => {
    const isPublic = publicRoutes.some(route => {
        const regex = new RegExp(`^${route.path.replace(/:\w+/g, '[^/]+')}$`);
        return regex.test(req.path) && route.method === req.method;
    });

    if (isPublic) {
        return next();
    }

    const role = req.user ? req.user.role : roles.USER;
    const userPermissions = permissions[role] || [];

    if (userPermissions.includes('*')) {
        return next(); // Admins can do anything
    }

    const matchedPermission = userPermissions.find(p => {
        const regex = new RegExp(`^${p.path.replace(/:\w+/g, '[^/]+')}$`);
        return regex.test(req.path) && p.method === req.method;
    });

    if (matchedPermission) {
        if (matchedPermission.own === true) {
            const requesterId = req.user.user_id.toString();
            // A more robust way to extract the user ID from the path
            const resourceId = req.params.userId || (req.path.split('/')[3]);

            if (requesterId === resourceId) {
                return next();
            }
            return res.status(403).json({ error: 'Forbidden: You can only access your own resources.' });
        } else if (matchedPermission.own === false) {
            // This permission is for moderators/admins to access any user's resource
            return next();
        }
        return next();
    }

    return res.status(403).json({ error: 'Forbidden: You do not have permission to access this resource.' });
};

module.exports = rbacMiddleware;
