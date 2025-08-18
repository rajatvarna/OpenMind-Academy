/**
 * This file defines the Role-Based Access Control (RBAC) system for the API Gateway.
 */

// Defines the available roles in the system.
const roles = {
  USER: 'user',
  MODERATOR: 'moderator',
  ADMIN: 'admin',
};

/**
 * Defines the permissions for each role.
 * The key is the role name, and the value is an array of permission objects.
 * A permission object can be a string ('*') for universal access, or an object with:
 *   - path: The URL path, which can include parameters like :userId.
 *   - method: The HTTP method (e.g., 'GET', 'POST').
 *   - own: (Optional) A boolean. If true, the user can only access the resource if they are the owner.
 *   - param: (Optional) If `own` is true, this specifies the URL parameter to use for the ownership check.
 */
const permissions = {
  [roles.ADMIN]: ['*'], // Admin can do anything
  [roles.MODERATOR]: [
    // Moderators can view other users' profiles and progress.
    { path: '/api/users/:userId/progress', method: 'GET', own: false },
    { path: '/api/users/:userId/full-profile', method: 'GET', own: false },
    { path: '/api/gamification/users/:userId/stats', method: 'GET', own: false },
    // General content creation permissions
    { path: '/api/content/reviews', method: 'POST' },
    { path: '/api/ugc/submit', method: 'POST' },
    { path: '/api/ugc/report', method: 'POST' },
    { path: '/api/qna/query', method: 'POST' },
    { path: '/api/qna/generate-quiz', method: 'POST' },
  ],
  [roles.USER]: [
    // Users can access their own profile and progress.
    { path: '/api/users/profile', method: 'GET', own: true }, // No param needed, tied to the user's own token.
    { path: '/api/users/:userId/progress', method: 'GET', own: true, param: 'userId' },
    { path: '/api/users/:userId/progress', method: 'POST', own: true, param: 'userId' },
    { path: '/api/users/:userId/full-profile', method: 'GET', own: true, param: 'userId' },
    { path: '/api/gamification/users/:userId/stats', method: 'GET', own: true, param: 'userId' },
    // General content creation permissions
    { path: '/api/content/reviews', method: 'POST' },
    { path: '/api/ugc/submit', method: 'POST' },
    { path: '/api/ugc/report', method: 'POST' },
    { path: '/api/qna/query', method: 'POST' },
    { path: '/api/qna/generate-quiz', method: 'POST' },
  ],
};

/**
 * A list of routes that are publicly accessible and do not require authentication.
 */
const publicRoutes = [
  { path: '/api/users/register', method: 'POST' },
  { path: '/api/users/login', method: 'POST' },
  { path: '/api/users/login/2fa', method: 'POST' },
  { path: '/api/content/courses', method: 'GET' },
  { path: '/api/content/courses/featured', method: 'GET' },
  { path: '/api/content/courses/:courseId', method: 'GET' },
  { path: '/api/content/courses/:courseId/reviews', method: 'GET' },
  { path: '/api/gamification/leaderboard', method: 'GET' },
  { path: '/health', method: 'GET' },
];


module.exports = {
  roles,
  permissions,
  publicRoutes,
};
