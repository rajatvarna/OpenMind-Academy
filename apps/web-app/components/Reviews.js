import { useState, useEffect } from 'react';
import styles from '../styles/Reviews.module.css';

// Star rating component
const StarRating = ({ rating }) => {
  return (
    <div className={styles.starRating}>
      {[...Array(5)].map((_, index) => (
        <span key={index} className={index < rating ? styles.filled : ''}>★</span>
      ))}
    </div>
  );
};

export default function Reviews({ courseId }) {
  const [reviews, setReviews] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!courseId) return;

    const fetchReviews = async () => {
      try {
        setIsLoading(true);
        // This will be a new API route we need to create
        const res = await fetch(`/api/courses/${courseId}/reviews`);
        if (res.ok) {
          const data = await res.json();
          setReviews(data);
        }
      } catch (error) {
        console.error('Failed to fetch reviews', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchReviews();
  }, [courseId]);

  if (isLoading) {
    return <div>Loading reviews...</div>;
  }

  if (reviews.length === 0) {
    return <div>No reviews yet. Be the first to leave one!</div>;
  }

  return (
    <div className={styles.reviewsContainer}>
      {reviews.map((review) => (
        <div key={review.id} className={styles.reviewCard}>
          <div className={styles.reviewHeader}>
            <strong>User {review.user_id}</strong>
            <StarRating rating={review.rating} />
          </div>
          <p>{review.review}</p>
          <small>{new Date(review.created_at).toLocaleDateString()}</small>
        </div>
      ))}
    </div>
  );
}
