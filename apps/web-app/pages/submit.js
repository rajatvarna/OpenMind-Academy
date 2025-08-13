import { useState, useEffect } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/Submit.module.css';

export default function SubmitPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  const [title, setTitle] = useState('');
  const [textContent, setTextContent] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    // If loading is finished and there's no user, redirect to login.
    if (!loading && !user) {
      router.push('/login');
    }
  }, [user, loading, router]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);

    try {
      // We need to create a lesson first in the Content Service, then submit to UGC service.
      // This is a multi-step process. For this example, we'll simplify and assume
      // we have a lessonId already, and we're just submitting the text for video generation.

      const lessonId = Date.now(); // Placeholder lessonId

      const res = await fetch('/api/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ lessonId, textContent, title }), // Sending title for context
      });

      if (res.ok) {
        alert('Submission successful! Your video is being generated.');
        router.push('/');
      } else {
        const data = await res.json();
        setError(data.message || 'Submission failed.');
      }
    } catch (err) {
      setError('An error occurred. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  // Render a loading state or null while checking auth
  if (loading || !user) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container">
      <Head>
        <title>Submit New Content</title>
      </Head>
      <main className={styles.main}>
        <h1 className={styles.title}>Submit New Content</h1>
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.inputGroup}>
            <label htmlFor="title">Title</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>
          <div className={styles.inputGroup}>
            <label htmlFor="textContent">Text for Video Generation</label>
            <textarea
              id="textContent"
              rows="10"
              value={textContent}
              onChange={(e) => setTextContent(e.target.value)}
              required
            />
          </div>
          <button type="submit" className={styles.button} disabled={submitting}>
            {submitting ? 'Submitting...' : 'Submit for Video Generation'}
          </button>
          {error && <p className={styles.error}>{error}</p>}
        </form>
      </main>
    </div>
  );
}
