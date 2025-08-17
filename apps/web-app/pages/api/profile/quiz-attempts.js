import { getCookie } from 'cookies-next';

const USER_SERVICE_URL = process.env.USER_SERVICE_URL || 'http://user-service:3000';

export default async function handler(req, res) {
  const token = getCookie('token', { req, res });

  if (!token) {
    return res.status(401).json({ message: 'Unauthorized' });
  }

  // We need to get the user ID from the token.
  // In a real app, you would decode the JWT here to get the user ID.
  // For now, we'll assume the client sends it. This is a simplification.
  // A better approach would be to have a /api/me endpoint that returns the user ID from the token.
  const { userId } = req.query;
  if (!userId) {
    return res.status(400).json({ message: 'User ID is required' });
  }


  if (req.method === 'GET') {
    try {
      const response = await fetch(`${USER_SERVICE_URL}/api/v1/users/${userId}/quiz-attempts`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        return res.status(response.status).json({ message: 'Failed to fetch quiz attempts' });
      }

      const data = await response.json();
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ message: 'Internal Server Error' });
    }
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
