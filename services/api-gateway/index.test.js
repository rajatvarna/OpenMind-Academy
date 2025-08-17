const jwt = require('jsonwebtoken');
const fs = require('fs');
const path = require('path');
const { generateKeyPairSync } = require('crypto');

describe('API Gateway Middlewares', () => {
  let mockRequest;
  let mockResponse;
  let authMiddleware;
  let rbacMiddleware;
  let nextFunction;
  let privateKey;
  const keyPath = path.join(__dirname, 'test_key.pub');

  beforeAll(() => {
    const { privateKey: genPrivateKey, publicKey } = generateKeyPairSync('rsa', {
      modulusLength: 2048,
      publicKeyEncoding: { type: 'spki', format: 'pem' },
      privateKeyEncoding: { type: 'pkcs8', format: 'pem' },
    });
    privateKey = genPrivateKey;
    fs.writeFileSync(keyPath, publicKey);
    process.env.JWT_PUBLIC_KEY_PATH = keyPath;

    // Load modules after setting the environment variable
    authMiddleware = require('./middleware/auth');
    rbacMiddleware = require('./middleware/rbac');
  });

  afterAll(() => {
    fs.unlinkSync(keyPath);
  });

  beforeEach(() => {
    mockRequest = {
      headers: {},
      user: null,
      params: {},
      path: '',
      method: '',
    };
    mockResponse = {
      status: jest.fn(() => mockResponse),
      json: jest.fn(),
      send: jest.fn(),
    };
    nextFunction = jest.fn();
  });

  describe('authMiddleware', () => {
    it('should call next() for a public route', () => {
      mockRequest.path = '/api/users/login';
      mockRequest.method = 'POST';
      authMiddleware(mockRequest, mockResponse, nextFunction);
      expect(nextFunction).toHaveBeenCalled();
    });

    it('should return 401 if no token is provided for a protected route', () => {
      mockRequest.path = '/api/users/profile';
      mockRequest.method = 'GET';
      authMiddleware(mockRequest, mockResponse, nextFunction);
      expect(mockResponse.status).toHaveBeenCalledWith(401);
    });

    it('should return 401 for an invalid token', () => {
      mockRequest.path = '/api/users/profile';
      mockRequest.method = 'GET';
      mockRequest.headers.authorization = 'Bearer invalid-token';
      authMiddleware(mockRequest, mockResponse, nextFunction);
      expect(mockResponse.status).toHaveBeenCalledWith(401);
    });

    it('should call next() and attach user for a valid token', () => {
        const user = { user_id: 1, role: 'user' };
        const token = jwt.sign(user, privateKey, { algorithm: 'RS256' });
        mockRequest.path = '/api/users/profile';
        mockRequest.method = 'GET';
        mockRequest.headers.authorization = `Bearer ${token}`;

        authMiddleware(mockRequest, mockResponse, nextFunction);

        expect(nextFunction).toHaveBeenCalled();
        expect(mockRequest.user).toBeDefined();
        expect(mockRequest.user.user_id).toBe(1);
    });
  });

  describe('rbacMiddleware', () => {
    it('should call next() for an admin', () => {
      mockRequest.user = { role: 'admin' };
      rbacMiddleware(mockRequest, mockResponse, nextFunction);
      expect(nextFunction).toHaveBeenCalled();
    });

    it('should call next() for a user with correct permissions', () => {
      mockRequest.user = { user_id: 1, role: 'user' };
      mockRequest.path = '/api/users/1/progress';
      mockRequest.method = 'GET';
      mockRequest.params.userId = '1';
      rbacMiddleware(mockRequest, mockResponse, nextFunction);
      expect(nextFunction).toHaveBeenCalled();
    });

    it('should return 403 for a user trying to access another user\'s resource', () => {
        mockRequest.user = { user_id: 1, role: 'user' };
        mockRequest.path = '/api/users/2/progress';
        mockRequest.method = 'GET';
        mockRequest.params.userId = '2';
        rbacMiddleware(mockRequest, mockResponse, nextFunction);
        expect(mockResponse.status).toHaveBeenCalledWith(403);
    });
  });
});
