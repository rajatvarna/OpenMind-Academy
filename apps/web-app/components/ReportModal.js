import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/ReportModal.module.css';

export default function ReportModal({ contentId, show, onClose }) {
  const { user } = useAuth();
  const [reason, setReason] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      await fetch('/api/report', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ contentId, reason, userId: user.id }),
      });
      alert('Report submitted. Thank you.');
      onClose();
    } catch (err) {
      alert('Failed to submit report.');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!show) {
    return null;
  }

  return (
    <div className={styles.modalBackdrop}>
      <div className={styles.modalContent}>
        <h2>Report Content</h2>
        <form onSubmit={handleSubmit}>
          <textarea
            rows="5"
            placeholder="Please provide a reason for your report..."
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            required
          />
          <div className={styles.buttons}>
            <button type="button" onClick={onClose} disabled={isSubmitting}>Cancel</button>
            <button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Submitting...' : 'Submit Report'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
