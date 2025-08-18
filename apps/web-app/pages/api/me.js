import cookie from 'cookie';

export default async function handler(req, res) {
  if (req.method !== 'GET') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;

  if (!token) {
    return res.status(401).json({ message: 'Not authenticated' });
  }

  try {
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const apiRes = await fetch(`${gatewayUrl}/api/users/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!apiRes.ok) {
      const errorData = await apiRes.json();
      return res.status(apiRes.status).json({ message: errorData.error || 'Failed to fetch user.' });
    }

    const user = await apiRes.json();
    res.status(200).json({ user });

  } catch (error) {
    console.error('Me API route error:', error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
