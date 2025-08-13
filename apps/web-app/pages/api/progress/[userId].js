import cookie from 'cookie';

export default async function handler(req, res) {
  const { userId } = req.query;

  // 1. Check for authentication
  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;
  if (!token) {
    return res.status(401).json({ message: 'Not authenticated.' });
  }
  // In a real app, you'd verify the token and ensure the token's userId matches the requested userId.

  const backendUrl = `http://api-gateway:8080/api/users/${userId}/progress`;

  if (req.method === 'GET') {
    try {
      const apiRes = await fetch(backendUrl, {
        headers: { 'Authorization': `Bearer ${token}` },
      });
      const data = await apiRes.json();
      if (!apiRes.ok) throw new Error(data.error || 'Failed to fetch progress');
      res.status(200).json(data);
    } catch (error) {
      res.status(500).json({ message: error.message });
    }
  } else if (req.method === 'POST') {
    try {
      const { lessonId } = req.body;
      const apiRes = await fetch(backendUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ lesson_id: lessonId }),
      });
      if (!apiRes.ok) {
        const errorData = await apiRes.json();
        throw new Error(errorData.error || 'Failed to update progress');
      }
      res.status(204).end();
    } catch (error) {
      res.status(500).json({ message: error.message });
    }
  } else {
    res.setHeader('Allow', ['GET', 'POST']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
