import React from 'react';
import Link from 'next/link';
import styles from '../styles/Header.module.css';

const Header = () => {
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
            <li>
              <Link href="/login" legacyBehavior>
                <a>Login</a>
              </Link>
            </li>
          </ul>
        </nav>
      </div>
    </header>
  );
};

export default Header;
