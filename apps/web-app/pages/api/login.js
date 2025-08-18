import { serialize } from 'cookie';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  const { email, password } = req.body;

  if (!email || !password) {
    return res.status(400).json({ message: 'Email and password are required.' });
  }

  try {
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    // Forward the login request to the actual User Service via the API Gateway
    const apiRes = await fetch(`${gatewayUrl}/api/users/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!apiRes.ok) {
      const errorData = await apiRes.json();
      // Forward the status code and error message from the backend
      return res.status(apiRes.status).json({ message: errorData.error || 'Authentication failed.' });
    }

    const data = await apiRes.json();

    // The backend will either return a full token, or a temp_token if 2FA is needed.
    // We just forward this response to the client.
    res.status(200).json(data);

  } catch (error) {
    console.error('Login API route error:', error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
