import { useState } from 'react';
import Modal from './Modal';
import styles from '../styles/DonationModal.module.css';

import { Elements, useStripe, useElements, PaymentElement } from '@stripe/react-stripe-js';

const CheckoutForm = ({ onCancel }) => {
  const stripe = useStripe();
  const elements = useElements();
  const [isProcessing, setIsProcessing] = useState(false);
  const [message, setMessage] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!stripe || !elements) {
      // Stripe.js has not yet loaded.
      // Make sure to disable form submission until Stripe.js has loaded.
      return;
    }

    setIsProcessing(true);

    const { error } = await stripe.confirmPayment({
      elements,
      confirmParams: {
        // Make sure to change this to your payment completion page
        return_url: `${window.location.origin}/`,
      },
    });

    if (error.type === "card_error" || error.type === "validation_error") {
      setMessage(error.message);
    } else {
      setMessage("An unexpected error occurred.");
    }

    setIsProcessing(false);
  };

  return (
    <form onSubmit={handleSubmit}>
      <PaymentElement />
      <div className={styles.buttonGroup}>
        <button type="button" onClick={onCancel} className={styles.cancelButton}>
          Cancel
        </button>
        <button disabled={isProcessing || !stripe || !elements} className={styles.submitButton}>
          {isProcessing ? "Processing..." : "Donate"}
        </button>
      </div>
      {/* Show any error or success messages */}
      {message && <div id="payment-message">{message}</div>}
    </form>
  );
};

export default function DonationModal({ show, onClose }) {
  const [amount, setAmount] = useState(500); // Default to $5.00
  const [clientSecret, setClientSecret] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleAmountSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const res = await fetch('/api/donations', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount }),
      });
      const data = await res.json();
      setClientSecret(data.clientSecret);
    } catch (error) {
      console.error("Failed to create payment intent", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal show={show} onClose={onClose} title="Support Free Education">
      {!clientSecret ? (
        <form onSubmit={handleAmountSubmit}>
          <p>Your support helps us keep education free and accessible for everyone.</p>
          <div className={styles.amountSelector}>
            <input
              type="number"
              value={amount / 100}
              onChange={(e) => setAmount(Number(e.target.value) * 100)}
              min="1"
            />
            <span>USD</span>
          </div>
          <button type="submit" disabled={isLoading} className={styles.submitButton}>
            {isLoading ? "Loading..." : "Proceed to Payment"}
          </button>
        </form>
      ) : (
        <Elements options={{ clientSecret }}>
          <CheckoutForm onCancel={() => setClientSecret(null)} />
        </Elements>
      )}
    </Modal>
  );
}
