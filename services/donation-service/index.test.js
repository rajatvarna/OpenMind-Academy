const request = require('supertest');
const app = require('./index');
const stripe = require('stripe');

jest.mock('stripe', () => {
  const originalStripe = jest.requireActual('stripe');
  return jest.fn((...args) => {
    const stripeInstance = originalStripe(...args);
    stripeInstance.paymentIntents = {
      create: jest.fn().mockResolvedValue({ client_secret: 'test_secret' }),
    };
    return stripeInstance;
  });
});

describe('Donation Service API', () => {
  let server;
  beforeAll((done) => {
    server = app.listen(0, done); // Listen on a random port
  });

  afterAll((done) => {
    server.close(done);
  });

  describe('POST /create-payment-intent', () => {
    it('should create a payment intent and return a client secret', async () => {
      const response = await request(app)
        .post('/create-payment-intent')
        .send({ amount: 1000 });

      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('clientSecret', 'test_secret');
      expect(stripe().paymentIntents.create).toHaveBeenCalledWith({
        amount: 1000,
        currency: 'usd',
        automatic_payment_methods: {
          enabled: true,
        },
      });
    });

    it('should return a 400 error if the amount is missing', async () => {
      const response = await request(app)
        .post('/create-payment-intent')
        .send({});

      expect(response.status).toBe(400);
      expect(response.body).toHaveProperty('error', 'Invalid amount.');
    });

    it('should return a 400 error if the amount is less than 100', async () => {
      const response = await request(app)
        .post('/create-payment-intent')
        .send({ amount: 50 });

      expect(response.status).toBe(400);
      expect(response.body).toHaveProperty('error', 'Invalid amount.');
    });

    it('should return a 500 error if Stripe fails', async () => {
      stripe().paymentIntents.create.mockRejectedValueOnce(new Error('Stripe error'));

      const response = await request(app)
        .post('/create-payment-intent')
        .send({ amount: 1000 });

      expect(response.status).toBe(500);
      expect(response.body).toHaveProperty('error', 'Failed to create payment intent.');
    });
  });
});
