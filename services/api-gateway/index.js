const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const helmet = require('helmet');
const cors = require('cors');
const morgan = require('morgan');
const jwt = require('jsonwebtoken');
const fs = require('fs');
require('dotenv').config();

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 8080;

// Middleware
app.use(helmet()); // Basic security headers
app.use(cors());   // Enable Cross-Origin Resource Sharing
app.use(morgan('combined')); // Request logging
app.use(express.json()); // To parse JSON bodies

// --- Authentication Middleware ---
const publicRoutes = [
    '/api/users/login',
    '/api/users/register',
    '/health'
];

const publicKey = fs.readFileSync('../secrets/jwtRS256.key.pub');

const authMiddleware = (req, res, next) => {
    if (publicRoutes.some(path => req.path.startsWith(path))) {
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

// --- Service Routes ---

// In a Kubernetes environment, 'user-service', 'content-service', etc.,
// would be the names of the Kubernetes services. The gateway would resolve
// these names to the correct internal IP addresses.

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

// Set up the proxy for each service
services.forEach(({ route, target }) => {
    app.use(route, createProxyMiddleware({
        target,
        changeOrigin: true,
        pathRewrite: {
            [`^${route}`]: '', // remove base path
        },
        onProxyReq: (proxyReq, req, res) => {
            // Add user identity headers to the downstream request
            if (req.user) {
                proxyReq.setHeader('X-User-Id', req.user.user_id);
                proxyReq.setHeader('X-User-Role', req.user.role);
            }
            console.log(`Proxying request for user ${req.user ? req.user.user_id : 'Guest'} to: ${target}${req.originalUrl}`);
        },
        onError: (err, req, res) => {
            console.error('Proxy error:', err);
            res.status(500).send('Proxy error');
        }
    }));
});

// Health check endpoint
app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

// Start the server
app.listen(PORT, () => {
    console.log(`API Gateway listening on port ${PORT}`);
});
