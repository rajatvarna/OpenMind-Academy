import { useState } from 'react';
import Modal from './Modal';
import styles from '../styles/DonationModal.module.css';

// This is a placeholder for the Stripe Elements wrapper
const CheckoutForm = ({ clientSecret }) => {
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsProcessing(true);
    // In a real app, you would use stripe.confirmPayment here
    console.log("Simulating payment confirmation with client secret:", clientSecret);
    setTimeout(() => {
      alert("Thank you for your donation!");
      setIsProcessing(false);
    }, 2000);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className={styles.placeholderCard}>
        <p>This is where the secure Stripe Card Element would be.</p>
        <span>**** **** **** 4242</span>
      </div>
      <button disabled={isProcessing} className={styles.submitButton}>
        {isProcessing ? "Processing..." : "Donate"}
      </button>
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
        // In a real app, you would wrap CheckoutForm with <Elements stripe={stripePromise} options={{ clientSecret }}>
        <CheckoutForm clientSecret={clientSecret} />
      )}
    </Modal>
  );
}
