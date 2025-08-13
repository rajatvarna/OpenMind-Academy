import { serialize } from 'cookie';

export default async function handler(req, res) {
  if (req.method !== 'POST' && req.method !== 'GET') { // Allow GET for simple link-based logout
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  // Set the cookie to be expired
  const cookie = serialize('auth_token', '', {
    maxAge: -1, // Expire the cookie
    path: '/',
  });

  res.setHeader('Set-Cookie', cookie);
  res.status(200).json({ success: true, message: 'Logged out successfully' });
}
