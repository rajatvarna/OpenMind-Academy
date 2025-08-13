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
  const [nextCursor, setNextCursor] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);

  const fetchReviews = async (cursor) => {
    try {
      const res = await fetch(`/api/courses/${courseId}/reviews?cursor=${cursor}&limit=5`);
      if (res.ok) {
        const data = await res.json();
        // Append new reviews to the existing list
        setReviews(prev => cursor === 0 ? data.data : [...prev, ...data.data]);
        setNextCursor(data.next_cursor);
      }
    } catch (error) {
      console.error('Failed to fetch reviews', error);
    }
  };

  useEffect(() => {
    if (!courseId) return;
    setIsLoading(true);
    fetchReviews(0).finally(() => setIsLoading(false));
  }, [courseId]);

  const handleLoadMore = () => {
    if (!nextCursor) return;
    setIsLoadingMore(true);
    fetchReviews(nextCursor).finally(() => setIsLoadingMore(false));
  };

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
      {nextCursor > 0 && (
        <button onClick={handleLoadMore} disabled={isLoadingMore} className={styles.loadMoreButton}>
          {isLoadingMore ? 'Loading...' : 'Load More Reviews'}
        </button>
      )}
    </div>
  );
}
