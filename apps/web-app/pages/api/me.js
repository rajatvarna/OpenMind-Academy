import cookie from 'cookie';
import { verify } from 'jsonwebtoken'; // In a real app, you'd use a JWT library

// A placeholder secret key. In a real app, this MUST be the same secret
// used by the User Service and stored securely as an environment variable.
const JWT_SECRET = process.env.JWT_SECRET || 'a-very-insecure-default-secret-key';

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
    // In a real application, you might not decode the token here. Instead, you might
    // forward it to the User Service to get the full, up-to-date user profile.
    // This prevents the user data in the token from becoming stale.

    // For this example, we'll simulate decoding it.
    // Note: The `verify` function would throw an error if the token is invalid or expired.
    // const decoded = verify(token, JWT_SECRET);

    // Since we can't actually verify without the real secret/library setup,
    // we'll just assume it's valid and return a dummy user.
    const dummyUser = { id: 1, email: 'user@example.com', name: 'Test User' };

    res.status(200).json({ user: dummyUser });
  } catch (error) {
    console.error('Me API route error:', error);
    res.status(401).json({ message: 'Invalid token' });
  }
}
