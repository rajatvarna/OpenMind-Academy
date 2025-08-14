const request = require('supertest');
const app = require('./index');
const stripe = require('stripe');

jest.mock('stripe', () => {
  const mockCreate = jest.fn().mockResolvedValue({ client_secret: 'test_secret' });
  return jest.fn().mockImplementation(() => ({
    paymentIntents: {
      create: mockCreate,
    },
  }));
});

describe('Donation Service API', () => {
  beforeEach(() => {
    const stripe = require('stripe');
    stripe().paymentIntents.create.mockClear();
    stripe().paymentIntents.create.mockResolvedValue({ client_secret: 'test_secret' });
  });

  describe('POST /create-payment-intent', () => {
    it('should create a payment intent and return a client secret', async () => {
      const stripe = require('stripe');
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
      const stripe = require('stripe');
      stripe().paymentIntents.create.mockRejectedValueOnce(new Error('Stripe error'));

      const response = await request(app)
        .post('/create-payment-intent')
        .send({ amount: 1000 });

      expect(response.status).toBe(500);
      expect(response.body).toHaveProperty('error', 'Failed to create payment intent.');
    });
  });
});
