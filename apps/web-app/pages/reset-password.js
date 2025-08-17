import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import styles from '../styles/Login.module.css'; // Reuse login styles

export default function ResetPasswordPage() {
  const router = useRouter();
  const [token, setToken] = useState(null);
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (router.isReady) {
      setToken(router.query.token);
    }
  }, [router.isReady, router.query.token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setMessage('');

    try {
      const res = await fetch('/api/password/reset', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, new_password: password }),
      });

      const data = await res.json();

      if (res.ok) {
        setMessage(data.message);
        setTimeout(() => router.push('/login'), 3000); // Redirect to login after 3 seconds
      } else {
        setError(data.message || 'An error occurred.');
      }
    } catch (err) {
      setError('An unexpected error occurred.');
    } finally {
      setLoading(false);
    }
  };

  if (!token) {
    return <div className={styles.container}>Loading...</div>;
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Reset Password</h1>
      <form onSubmit={handleSubmit} className={styles.form}>
        <input
          type="password"
          placeholder="New Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength="8"
          className={styles.input}
        />
        <button type="submit" disabled={loading} className={styles.button}>
          {loading ? 'Resetting...' : 'Reset Password'}
        </button>
        {message && <p className={styles.message}>{message}</p>}
        {error && <p className={styles.error}>{error}</p>}
      </form>
    </div>
  );
}
