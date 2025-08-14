export default async function handler(req, res) {
  if (req.method !== 'GET') {
    res.setHeader('Allow', ['GET']);
    return res.status(405).end(`Method ${req.method} Not Allowed`);
  }

  try {
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const backendUrl = `${gatewayUrl}/api/content/courses/featured`;
    const apiRes = await fetch(backendUrl);

    if (!apiRes.ok) {
      const errorData = await apiRes.json();
      return res.status(apiRes.status).json({ message: errorData.error || 'Failed to fetch featured courses' });
    }

    const courses = await apiRes.json();
    res.status(200).json(courses);

  } catch (error) {
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
