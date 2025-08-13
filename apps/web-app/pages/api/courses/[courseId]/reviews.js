export default async function handler(req, res) {
  const { courseId } = req.query;

  if (req.method !== 'GET') {
    res.setHeader('Allow', ['GET']);
    return res.status(405).end(`Method ${req.method} Not Allowed`);
  }

  try {
    const apiRes = await fetch(`http://api-gateway:8080/api/content/courses/${courseId}/reviews`);

    if (!apiRes.ok) {
      const errorData = await apiRes.json();
      return res.status(apiRes.status).json({ message: errorData.error || 'Failed to fetch reviews' });
    }

    const reviews = await apiRes.json();
    res.status(200).json(reviews);

  } catch (error) {
    console.error(`Failed to fetch reviews for course ${courseId}:`, error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
