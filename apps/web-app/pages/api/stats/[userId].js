import cookie from 'cookie';

export default async function handler(req, res) {
  const { userId } = req.query;

  // Check for authentication
  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;
  if (!token) {
    return res.status(401).json({ message: 'Not authenticated.' });
  }
  // You would also verify the token and that the user is requesting their own stats.

  if (req.method === 'GET') {
    try {
      const backendUrl = `http://api-gateway:8080/api/gamification/users/${userId}/stats`;
      // Assuming the gamification service is mounted at /api/gamification on the gateway

      const apiRes = await fetch(backendUrl, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      const data = await apiRes.json();
      if (!apiRes.ok) throw new Error(data.error || 'Failed to fetch stats');

      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ message: error.message });
    }
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
