import React, { useState, useEffect } from 'react';
import './ModerationQueue.css';

// --- Placeholder Data ---
// In a real app, this would be fetched from an API endpoint like GET /api/ugc/pending
const placeholderSubmissions = [
  { id: 101, title: 'Introduction to Quantum Physics', author: 'User123', status: 'pending' },
  { id: 102, title: 'History of the Roman Empire', author: 'HistoryBuff', status: 'pending' },
  { id: 103, title: 'Advanced Baking Techniques', author: 'BakerPro', status: 'pending' },
  { id: 104, title: 'How to Build a REST API in Go', author: 'GoDev', status: 'pending' },
];

function ModerationQueue({ onSelectSubmission }) {
  const [submissions, setSubmissions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Simulate fetching data from an API
    const fetchSubmissions = () => {
      try {
        // In a real app:
        // const response = await fetch('/api/ugc/pending');
        // const data = await response.json();
        // setSubmissions(data);
        setSubmissions(placeholderSubmissions);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch submissions.');
        setLoading(false);
      }
    };

    fetchSubmissions();
  }, []); // The empty dependency array ensures this effect runs only once on mount

  if (loading) {
    return <div>Loading queue...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="moderation-queue">
      <h2>Pending Submissions</h2>
      <ul>
        {submissions.map((submission) => (
          <li key={submission.id} onClick={() => onSelectSubmission(submission)}>
            <div className="submission-title">{submission.title}</div>
            <div className="submission-author">by {submission.author}</div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ModerationQueue;
