const { roles, permissions, publicRoutes } = require('./rbac_config');

// This is a simplified version of the rbacMiddleware from index.js for testing
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


describe('RBAC Middleware', () => {
  let mockRequest;
  let mockResponse;
  let nextFunction;

  beforeEach(() => {
    mockRequest = {
      headers: {},
      user: null,
      params: {},
    };
    mockResponse = {
      status: jest.fn(() => mockResponse),
      json: jest.fn(),
    };
    nextFunction = jest.fn();
  });

  // Test cases for admin
  test('Admin should be able to access any resource', () => {
    mockRequest.user = { role: 'admin' };
    mockRequest.path = '/api/users/123/full-profile';
    mockRequest.method = 'GET';
    rbacMiddleware(mockRequest, mockResponse, nextFunction);
    expect(nextFunction).toHaveBeenCalled();
  });

  // Test cases for moderator
  test('Moderator should be able to access other users resources', () => {
    mockRequest.user = { role: 'moderator' };
    mockRequest.path = '/api/users/456/full-profile';
    mockRequest.method = 'GET';
    rbacMiddleware(mockRequest, mockResponse, nextFunction);
    expect(nextFunction).toHaveBeenCalled();
  });

  // Test cases for user
  test('User should be able to access their own resources', () => {
    mockRequest.user = { user_id: '123', role: 'user' };
    mockRequest.path = '/api/users/123/full-profile';
    mockRequest.method = 'GET';
    rbacMiddleware(mockRequest, mockResponse, nextFunction);
    expect(nextFunction).toHaveBeenCalled();
  });

  test('User should be denied access to other users resources', () => {
    mockRequest.user = { user_id: '123', role: 'user' };
    mockRequest.path = '/api/users/456/full-profile';
    mockRequest.method = 'GET';
    rbacMiddleware(mockRequest, mockResponse, nextFunction);
    expect(mockResponse.status).toHaveBeenCalledWith(403);
  });

  // Test cases for public routes
  test('Public route should be accessible without authentication', () => {
    mockRequest.path = '/api/content/courses';
    mockRequest.method = 'GET';
    rbacMiddleware(mockRequest, mockResponse, nextFunction);
    expect(nextFunction).toHaveBeenCalled();
  });
});
