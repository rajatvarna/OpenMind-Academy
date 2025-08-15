const roles = {
  USER: 'user',
  MODERATOR: 'moderator',
  ADMIN: 'admin',
};

const permissions = {
    [roles.ADMIN]: ['*'], // Admin can do anything
    [roles.MODERATOR]: [
        { path: '/api/users/profile', method: 'GET' },
        { path: '/api/users/:userId/progress', method: 'GET' },
        { path: '/api/users/:userId/full-profile', method: 'GET' },
        { path: '/api/content/reviews', method: 'POST' },
        { path: '/api/ugc/submit', method: 'POST' },
        { path: '/api/ugc/report', method: 'POST' },
        { path: '/api/qna/query', method: 'POST' },
        { path: '/api/qna/generate-quiz', method: 'POST' },
        { path: '/api/gamification/users/:userId/stats', method: 'GET' },
    ],
    [roles.USER]: [
        { path: '/api/users/profile', method: 'GET', own: true },
        { path: '/api/users/:userId/progress', method: 'GET', own: true },
        { path: '/api/users/:userId/progress', method: 'POST', own: true },
        { path: '/api/users/:userId/full-profile', method: 'GET', own: true },
        { path: '/api/content/reviews', method: 'POST' },
        { path: '/api/ugc/submit', method: 'POST' },
        { path: '/api/ugc/report', method: 'POST' },
        { path: '/api/qna/query', method: 'POST' },
        { path: '/api/qna/generate-quiz', method: 'POST' },
        { path: '/api/gamification/users/:userId/stats', method: 'GET', own: true },
    ],
};

module.exports = {
  roles,
  permissions,
};
