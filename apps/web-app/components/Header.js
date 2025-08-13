import React from 'react';
import Link from 'next/link';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/Header.module.css';

const Header = () => {
  const { user, logout, loading } = useAuth();

  return (
    <header className={styles.header}>
      <div className={styles.container}>
        <Link href="/" legacyBehavior>
          <a className={styles.logo}>FreeEdu</a>
        </Link>
        <nav>
          <ul className={styles.navList}>
            <li>
              <Link href="/" legacyBehavior>
                <a>Courses</a>
              </Link>
            </li>
            {loading ? null : user ? (
              <>
                <li className={styles.userName}>Hello, {user.name}</li>
                <li>
                  <button onClick={logout} className={styles.logoutButton}>Logout</button>
                </li>
              </>
            ) : (
              <li>
                <Link href="/login" legacyBehavior>
                  <a>Login</a>
                </Link>
              </li>
            )}
          </ul>
        </nav>
      </div>
    </header>
  );
};

export default Header;
