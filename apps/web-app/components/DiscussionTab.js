import { useState, useEffect } from 'react';
import styles from '../styles/DiscussionTab.module.css';

// This is a simplified discussion component. A real one would be more complex.
export default function DiscussionTab({ courseId }) {
  const [threads, setThreads] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!courseId) return;
    const fetchThreads = async () => {
      try {
        setIsLoading(true);
        // This will be a new API route
        const res = await fetch(`/api/courses/${courseId}/threads`);
        if (res.ok) {
          const data = await res.json();
          setThreads(data);
        }
      } catch (error) {
        console.error('Failed to fetch threads', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchThreads();
  }, [courseId]);

  if (isLoading) return <p>Loading discussions...</p>;

  return (
    <div className={styles.discussionContainer}>
      <div className={styles.newThread}>
        {/* We would add a form here to create a new thread */}
        <button>Start a New Discussion</button>
      </div>
      <div className={styles.threadList}>
        {threads.length > 0 ? (
          threads.map(thread => (
            <div key={thread.id} className={styles.threadItem}>
              <p className={styles.threadTitle}>{thread.title}</p>
              <span className={styles.threadMeta}>by User {thread.user_id} on {new Date(thread.created_at).toLocaleDateString()}</span>
            </div>
          ))
        ) : (
          <p>No discussions yet. Start one!</p>
        )}
      </div>
    </div>
  );
}
