import { getCookie } from 'cookies-next';

const USER_SERVICE_URL = process.env.USER_SERVICE_URL || 'http://user-service:3000';

export default async function handler(req, res) {
  const token = getCookie('token', { req, res });

  if (!token) {
    return res.status(401).json({ message: 'Unauthorized' });
  }

  if (req.method === 'GET') {
    try {
      const response = await fetch(`${USER_SERVICE_URL}/api/v1/preferences`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        return res.status(response.status).json({ message: 'Failed to fetch preferences' });
      }

      const data = await response.json();
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ message: 'Internal Server Error' });
    }
  } else if (req.method === 'PUT') {
    try {
      const response = await fetch(`${USER_SERVICE_URL}/api/v1/preferences`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(req.body),
      });

      if (!response.ok) {
        return res.status(response.status).json({ message: 'Failed to update preferences' });
      }

      res.status(204).end();
    } catch (error) {
      res.status(500).json({ message: 'Internal Server Error' });
    }
  } else {
    res.setHeader('Allow', ['GET', 'PUT']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
