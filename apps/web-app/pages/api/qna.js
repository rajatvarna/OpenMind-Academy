export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  const { question } = req.body;

  if (!question) {
    return res.status(400).json({ message: 'A question is required.' });
  }

  try {
    // Forward the request to the AI Q&A Service via the API Gateway
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const qnaApiRes = await fetch(`${gatewayUrl}/api/qna/query`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ question }),
    });

    const data = await qnaApiRes.json();

    if (!qnaApiRes.ok) {
      return res.status(qnaApiRes.status).json({ message: data.detail || 'Failed to get an answer.' });
    }

    // Return the response from the Q&A service to the client
    res.status(200).json(data);

  } catch (error) {
    console.error('Q&A API route error:', error);
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
