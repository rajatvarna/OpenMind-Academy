import cookie from 'cookie';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;
  if (!token) {
    return res.status(401).json({ message: 'Not authenticated.' });
  }

  const { contentId, reason, userId } = req.body;
  if (!contentId || !reason || !userId) {
    return res.status(400).json({ message: 'Content ID, reason, and user ID are required.' });
  }

  try {
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const backendUrl = `${gatewayUrl}/api/ugc/report`;
    const apiRes = await fetch(backendUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ contentId, reason, userId }),
    });

    const data = await apiRes.json();
    if (!apiRes.ok) {
      return res.status(apiRes.status).json({ message: data.error || 'Failed to submit report' });
    }

    res.status(202).json(data);

  } catch (error) {
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
