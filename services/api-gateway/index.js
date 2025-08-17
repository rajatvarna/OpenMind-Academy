const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const helmet = require('helmet');
const cors = require('cors');
const morgan = require('morgan');
const jwt = require('jsonwebtoken');
const fs = require('fs');
require('dotenv').config();
const { roles, permissions, publicRoutes } = require('./rbac_config');

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 8080;

// --- Core Middleware ---
app.use(helmet()); // Set various HTTP headers for security
app.use(cors());   // Enable Cross-Origin Resource Sharing for all routes
app.use(morgan('combined')); // Log HTTP requests
app.use(express.json()); // Parse incoming JSON requests

// --- Authentication Middleware ---
// This middleware is responsible for verifying the JWT token.
// It checks for the Authorization header, verifies the token, and attaches the decoded user to the request.
// Public routes are skipped by this middleware.
const publicKey = fs.readFileSync('../secrets/jwtRS256.key.pub');

const authMiddleware = (req, res, next) => {
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

app.use(authMiddleware);

// --- RBAC (Role-Based Access Control) Middleware ---
// This middleware enforces permissions based on user roles.
// It checks the user's role (from the JWT token) against the permissions defined in `rbac_config.js`.
// It supports wildcard permissions for admins, and "own" resource checks for users.
const rbacMiddleware = (req, res, next) => {
    // Public routes are allowed to bypass the RBAC check.
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

app.use(rbacMiddleware);

// --- Service Routes ---

// --- Service-to-Route Mapping ---
// This configuration maps API routes to their corresponding backend microservices.
// The `target` URLs would typically be the internal DNS names of the services in a container orchestration environment (e.g., Kubernetes).
const services = [
    {
        route: '/api/users',
        target: process.env.USER_SERVICE_URL || 'http://user-service:3000',
    },
    {
        route: '/api/content',
        target: process.env.CONTENT_SERVICE_URL || 'http://content-service:3001',
    },
    {
        route: '/api/ugc',
        target: process.env.UGC_SERVICE_URL || 'http://ugc-service:3002',
    },
    {
        route: '/api/qna',
        target: process.env.QNA_SERVICE_URL || 'http://qna-service:3003',
    },
    // Add other services here as they are built
];

// --- Proxy Middleware Setup ---
// Dynamically create proxy middleware for each service defined above.
// The `http-proxy-middleware` library handles the request forwarding.
services.forEach(({ route, target }) => {
    // Options for the proxy middleware
    const proxyOptions = {
        target,
        changeOrigin: true, // Needed for virtual hosted sites
        pathRewrite: {
            [`^${route}`]: '', // Rewrite path: remove the base route segment
        },
        onProxyReq: (proxyReq, req, res) => {
            // Forward user identity to downstream services
            if (req.user) {
                proxyReq.setHeader('X-User-Id', req.user.user_id);
                proxyReq.setHeader('X-User-Role', req.user.role);
            }
            console.log(`Proxying request for user ${req.user ? req.user.user_id : 'Guest'} to: ${target}${req.originalUrl}`);
        },
        onError: (err, req, res) => {
            console.error('Proxy error:', err);
            res.status(500).send('Proxy Error');
        }
    };

    app.use(route, createProxyMiddleware(proxyOptions));
});

// --- Health Check Endpoint ---
// A simple endpoint to verify that the gateway is running.
// This is useful for load balancers and container orchestrators.
app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

// --- Server Activation ---
app.listen(PORT, () => {
    console.log(`API Gateway listening on port ${PORT}`);
});
