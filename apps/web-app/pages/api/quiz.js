export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  const { text_content } = req.body;
  if (!text_content) {
    return res.status(400).json({ message: 'Text content is required.' });
  }

  try {
    const backendUrl = 'http://api-gateway:8080/api/qna/generate-quiz';
    const apiRes = await fetch(backendUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text_content }),
    });

    const data = await apiRes.json();
    if (!apiRes.ok) {
      return res.status(apiRes.status).json({ message: data.detail || 'Failed to generate quiz' });
    }

    res.status(200).json(data);

  } catch (error) {
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
