export default async function handler(req, res) {
  if (req.method !== 'GET') {
    res.setHeader('Allow', ['GET']);
    return res.status(405).end(`Method ${req.method} Not Allowed`);
  }

  const { q } = req.query;
  if (!q) {
    return res.status(400).json({ message: 'Query parameter "q" is required.' });
  }

  try {
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const backendUrl = `${gatewayUrl}/api/search/search?q=${encodeURIComponent(q)}`;
    // Assuming search service is at /api/search on the gateway

    const apiRes = await fetch(backendUrl);

    if (!apiRes.ok) {
      const errorData = await apiRes.json();
      return res.status(apiRes.status).json({ message: errorData.error || 'Failed to fetch search results' });
    }

    const searchData = await apiRes.json();
    res.status(200).json(searchData);

  } catch (error) {
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
