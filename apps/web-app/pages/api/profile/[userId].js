import cookie from 'cookie';

export default async function handler(req, res) {
  const { userId } = req.query;

  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;
  if (!token) {
    return res.status(401).json({ message: 'Not authenticated.' });
  }

  if (req.method === 'GET') {
    try {
      const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
      const backendUrl = `${gatewayUrl}/api/users/${userId}/full-profile`;
      const apiRes = await fetch(backendUrl, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      const data = await apiRes.json();
      if (!apiRes.ok) throw new Error(data.error || 'Failed to fetch profile data');

      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ message: error.message });
    }
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
