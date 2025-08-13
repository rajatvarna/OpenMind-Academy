export default async function handler(req, res) {
  const { courseId } = req.query;

  if (req.method !== 'GET') {
    res.setHeader('Allow', ['GET']);
    return res.status(405).end(`Method ${req.method} Not Allowed`);
  }

  try {
    const backendUrl = `http://api-gateway:8080/api/forum/courses/${courseId}/threads`;
    // Assuming the forum service is mounted at /api/forum on the gateway

    const apiRes = await fetch(backendUrl);

    if (!apiRes.ok) {
      const errorData = await apiRes.json();
      return res.status(apiRes.status).json({ message: errorData.error || 'Failed to fetch threads' });
    }

    const threads = await apiRes.json();
    res.status(200).json(threads);

  } catch (error) {
    console.error(`Failed to fetch threads for course ${courseId}:`, error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
