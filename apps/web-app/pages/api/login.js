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

    const { token } = await apiRes.json();

    // Set the JWT in a secure, HttpOnly cookie.
    // This cookie will be automatically sent by the browser on subsequent requests.
    const cookie = serialize('auth_token', token, {
      httpOnly: true, // The cookie is not accessible via client-side JavaScript
      secure: process.env.NODE_ENV !== 'development', // Use secure cookies in production
      maxAge: 60 * 60 * 24 * 7, // 1 week
      sameSite: 'strict',
      path: '/',
    });

    res.setHeader('Set-Cookie', cookie);
    res.status(200).json({ success: true });

  } catch (error) {
    console.error('Login API route error:', error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
