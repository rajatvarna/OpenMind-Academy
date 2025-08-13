import cookie from 'cookie';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  // 1. Check for authentication
  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;
  if (!token) {
    return res.status(401).json({ message: 'Not authenticated.' });
  }

  // 2. Get data from request body
  const { courseId, userId, rating, review } = req.body;

  try {
    // 3. Forward to the Content Service
    const apiRes = await fetch('http://api-gateway:8080/api/content/reviews', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ course_id: courseId, user_id: userId, rating, review }),
    });

    const data = await apiRes.json();
    if (!apiRes.ok) {
      return res.status(apiRes.status).json({ message: data.error || 'Failed to submit review' });
    }

    res.status(201).json(data);

  } catch (error) {
    console.error('Submit review API error:', error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
