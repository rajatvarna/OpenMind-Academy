const jwt = require('jsonwebtoken');
const fs = require('fs');
const { publicRoutes } = require('../rbac_config');

const authMiddleware = (req, res, next) => {
    const publicKey = fs.readFileSync(process.env.JWT_PUBLIC_KEY_PATH || '../secrets/jwtRS256.key.pub');
    const isPublic = publicRoutes.some(route => {
        const regex = new RegExp(`^${route.path.replace(/:\w+/g, '[^/]+')}$`);
        return regex.test(req.path) && route.method === req.method;
    });

    if (isPublic) {
        return next();
    }

    const authHeader = req.headers.authorization;
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return res.status(401).json({ error: 'Authorization header is required' });
    }

    const token = authHeader.split(' ')[1];
    try {
        const decoded = jwt.verify(token, publicKey, { algorithms: ['RS256'] });
        req.user = decoded; // Attach decoded user info to the request
        next();
    } catch (err) {
        return res.status(401).json({ error: 'Invalid or expired token' });
    }
};

module.exports = authMiddleware;
