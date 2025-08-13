const express = require('express');
const cors = require('cors');
// In a real application, load this from an environment variable.
const stripe = require('stripe')(process.env.STRIPE_SECRET_KEY);

const app = express();
app.use(cors()); // Enable CORS for the frontend to call this
app.use(express.json());

app.post('/create-payment-intent', async (req, res) => {
  const { amount } = req.body; // Amount in cents

  if (!amount || amount < 100) { // Example: minimum donation of $1.00
    return res.status(400).send({ error: 'Invalid amount.' });
  }

  try {
    // Create a PaymentIntent with the order amount and currency
    const paymentIntent = await stripe.paymentIntents.create({
      amount: amount,
      currency: 'usd',
      automatic_payment_methods: {
        enabled: true,
      },
    });

    res.send({
      clientSecret: paymentIntent.client_secret,
    });
  } catch (error) {
    console.error("Stripe Error:", error.message);
    res.status(500).send({ error: 'Failed to create payment intent.' });
  }
});

const PORT = process.env.PORT || 3007; // Port for the donation service
app.listen(PORT, () => console.log(`Donation service running on port ${PORT}`));
