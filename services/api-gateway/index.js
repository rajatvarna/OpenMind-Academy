const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const helmet = require('helmet');
const cors = require('cors');
const morgan = require('morgan');

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 8080;

// Middleware
app.use(helmet()); // Basic security headers
app.use(cors());   // Enable Cross-Origin Resource Sharing
app.use(morgan('combined')); // Request logging
app.use(express.json()); // To parse JSON bodies

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
            // Here you could add custom headers, like a trace ID
            console.log(`Proxying request to: ${target}${req.originalUrl}`);
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
