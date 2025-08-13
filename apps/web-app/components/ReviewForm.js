import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/ReviewForm.module.css';

export default function ReviewForm({ courseId, onReviewSubmitted }) {
  const { user } = useAuth();
  const [rating, setRating] = useState(0);
  const [reviewText, setReviewText] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (rating === 0) {
      setError('Please select a rating.');
      return;
    }
    setError('');
    setIsSubmitting(true);

    try {
      const res = await fetch('/api/reviews', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          courseId,
          userId: user.id,
          rating,
          review: reviewText,
        }),
      });

      if (res.ok) {
        setRating(0);
        setReviewText('');
        if (onReviewSubmitted) {
          onReviewSubmitted(); // Notify parent to refetch reviews
        }
      } else {
        const data = await res.json();
        setError(data.message || 'Failed to submit review.');
      }
    } catch (err) {
      setError('An error occurred. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!user) {
    return <p>Please log in to leave a review.</p>;
  }

  return (
    <form onSubmit={handleSubmit} className={styles.reviewForm}>
      <h3>Write a Review</h3>
      <div className={styles.starRating}>
        {[...Array(5)].map((_, index) => {
          const ratingValue = index + 1;
          return (
            <button
              type="button"
              key={ratingValue}
              className={ratingValue <= rating ? styles.on : styles.off}
              onClick={() => setRating(ratingValue)}
            >
              <span className="star">★</span>
            </button>
          );
        })}
      </div>
      <textarea
        rows="4"
        value={reviewText}
        onChange={(e) => setReviewText(e.target.value)}
        placeholder="Share your thoughts on this course..."
      />
      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Submitting...' : 'Submit Review'}
      </button>
      {error && <p className={styles.error}>{error}</p>}
    </form>
  );
}
