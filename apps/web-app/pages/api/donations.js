export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }

  const { amount } = req.body;
  if (!amount) {
    return res.status(400).json({ message: 'Amount is required.' });
  }

  try {
    const backendUrl = 'http://api-gateway:8080/api/donations/create-payment-intent';
    // Assuming donation service is at /api/donations on the gateway

    const apiRes = await fetch(backendUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ amount }),
    });

    const data = await apiRes.json();
    if (!apiRes.ok) {
      return res.status(apiRes.status).json({ message: data.error || 'Failed to create payment intent' });
    }

    res.status(200).json(data);

  } catch (error) {
    res.status(500).json({ message: 'An internal server error occurred.' });
  }
}
