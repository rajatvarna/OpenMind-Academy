import cookie from 'cookie';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  // 1. Check for authentication
  const cookies = cookie.parse(req.headers.cookie || '');
  const token = cookies.auth_token;

  if (!token) {
    return res.status(401).json({ message: 'Not authenticated. Please log in.' });
  }

  // The allowlist check has been removed for public launch.
  // All authenticated users can now submit content.

  // 2. Get data from the request body
  const { lessonId, textContent, title } = req.body;

  if (!lessonId || !textContent || !title) {
    return res.status(400).json({ message: 'Lesson ID, text content, and title are required.' });
  }

  try {
    // 3. Forward the request to the UGC Submission Service
    const ugcApiRes = await fetch('http://api-gateway:8080/api/ugc/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // You could also forward the JWT for the backend service to use
        // 'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ lessonId, textContent }),
    });

    // 4. Return the response from the UGC service to the client
    const data = await ugcApiRes.json();
    if (!ugcApiRes.ok) {
      return res.status(ugcApiRes.status).json({ message: data.error || 'Submission failed.' });
    }

    res.status(202).json(data);

  } catch (error) {
    console.error('Submit API route error:', error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
