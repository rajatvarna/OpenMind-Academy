import { useState } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import styles from '../styles/Login.module.css';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();

  // State for 2FA flow
  const [needs2FA, setNeeds2FA] = useState(false);
  const [tempToken, setTempToken] = useState('');
  const [twoFactorToken, setTwoFactorToken] = useState('');

  const handlePasswordSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();
      if (res.ok) {
        if (data.temp_token) {
          // 2FA is needed
          setTempToken(data.temp_token);
          setNeeds2FA(true);
        } else {
          // Login was successful, the cookie is set by the API route.
          router.push('/');
        }
      } else {
        setError(data.message || 'Failed to login.');
      }
    } catch (err) {
      setError('An error occurred. Please try again.');
    }
  };

  const handle2FASubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const res = await fetch('/api/users/login/2fa', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ temp_token: tempToken, token: twoFactorToken }),
      });

      if (res.ok) {
        // 2FA login successful, cookie is set.
        router.push('/');
      } else {
        const data = await res.json();
        setError(data.error || 'Failed 2FA verification.');
      }
    } catch (err) {
      setError('An error occurred. Please try again.');
    }
  };

  return (
    <div className="container">
      <Head>
        <title>Login</title>
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>{needs2FA ? 'Enter 2FA Code' : 'Login'}</h1>

        {!needs2FA ? (
          <form onSubmit={handlePasswordSubmit} className={styles.form}>
            <div className={styles.inputGroup}>
              <label htmlFor="email">Email</label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className={styles.inputGroup}>
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <button type="submit" className={styles.button}>Login</button>
          </form>
        ) : (
          <form onSubmit={handle2FASubmit} className={styles.form}>
            <p>Enter the code from your authenticator app.</p>
            <div className={styles.inputGroup}>
              <label htmlFor="2fa-token">Authentication Code</label>
              <input
                type="text"
                id="2fa-token"
                value={twoFactorToken}
                onChange={(e) => setTwoFactorToken(e.target.value)}
                required
                maxLength="6"
              />
            </div>
            <button type="submit" className={styles.button}>Verify</button>
          </form>
        )}

        {error && <p className={styles.error}>{error}</p>}
      </main>
    </div>
  );
}
