const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const helmet = require('helmet');
const cors = require('cors');
const morgan = require('morgan');
require('dotenv').config();
const authMiddleware = require('./middleware/auth');
const rbacMiddleware = require('./middleware/rbac');

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 8080;

// --- Core Middleware ---
app.use(helmet()); // Set various HTTP headers for security
app.use(cors());   // Enable Cross-Origin Resource Sharing for all routes
app.use(morgan('combined')); // Log HTTP requests
app.use(express.json()); // Parse incoming JSON requests

app.use(authMiddleware);
app.use(rbacMiddleware);

// --- Service Routes ---
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
];

// --- Proxy Middleware Setup ---
services.forEach(({ route, target }) => {
    // Common proxy options
    const proxyOptions = {
        target,
        changeOrigin: true, // Necessary for virtual-hosted sites
        pathRewrite: {
            [`^${route}`]: '', // Rewrite path to remove the gateway-specific route prefix
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
            // Centralized error handling for proxy issues
            console.error('Proxy error:', err);
            res.status(500).send('Proxy Error');
        }
    };

    // Apply the proxy for the defined route
    app.use(route, createProxyMiddleware(proxyOptions));
});

// --- Health Check Endpoint ---
app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

// --- Server Activation ---
app.listen(PORT, () => {
    console.log(`API Gateway listening on port ${PORT}`);
});

module.exports = { app };
