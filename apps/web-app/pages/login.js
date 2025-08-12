import { useState } from 'react';
import Head from 'next/head';
import styles from '../styles/Login.module.css';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log('Logging in with:', { email, password });

    // In a real application, you would make an API call here:
    // const res = await fetch('/api/login', { // This might be a Next.js API route proxying to your user service
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ email, password }),
    // });
    // if (res.ok) {
    //   const { token } = await res.json();
    //   // Store the token (e.g., in cookies) and redirect
    //   console.log('Login successful, token:', token);
    // } else {
    //   console.error('Login failed');
    // }
    alert('Login functionality is not yet implemented. Check the console for details.');
  };

  return (
    <div className="container">
      <Head>
        <title>Login</title>
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>Login</h1>
        <form onSubmit={handleSubmit} className={styles.form}>
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
      </main>
    </div>
  );
}
