const roles = {
  USER: 'user',
  MODERATOR: 'moderator',
  ADMIN: 'admin',
};

const permissions = {
  [roles.ADMIN]: ['*'], // Admin can do anything
  [roles.MODERATOR]: [
    // User management
    { path: '/api/users/:userId/progress', method: 'GET', own: false },
    { path: '/api/users/:userId/full-profile', method: 'GET', own: false },
    { path: '/api/gamification/users/:userId/stats', method: 'GET', own: false },
    // Content management
    { path: '/api/content/reviews', method: 'POST' },
    { path: '/api/ugc/submit', method: 'POST' },
    { path: '/api/ugc/report', method: 'POST' },
    // Q&A
    { path: '/api/qna/query', method: 'POST' },
    { path: '/api/qna/generate-quiz', method: 'POST' },
  ],
  [roles.USER]: [
    // Own resource access
    { path: '/api/users/profile', method: 'GET', own: true },
    { path: '/api/users/:userId/progress', method: 'GET', own: true },
    { path: '/api/users/:userId/progress', method: 'POST', own: true },
    { path: '/api/users/:userId/full-profile', method: 'GET', own: true },
    { path: '/api/gamification/users/:userId/stats', method: 'GET', own: true },
    // General permissions
    { path: '/api/content/reviews', method: 'POST' },
    { path: '/api/ugc/submit', method: 'POST' },
    { path: '/api/ugc/report', method: 'POST' },
    { path: '/api/qna/query', method: 'POST' },
    { path: '/api/qna/generate-quiz', method: 'POST' },
  ],
};

const publicRoutes = [
  { path: '/api/users/register', method: 'POST' },
  { path: '/api/users/login', method: 'POST' },
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
