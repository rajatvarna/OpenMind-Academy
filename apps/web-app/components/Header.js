import React from 'react';
import Link from 'next/link';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/Header.module.css';
import SearchBar from './SearchBar';
import DonationModal from './DonationModal';
import { useState } from 'react';

/**
 * The main header for the web application.
 * It displays the site logo, navigation links, a search bar, and user-specific information.
 * It handles different states for logged-in and logged-out users.
 */
const Header = () => {
  // useAuth provides user session, stats, and authentication functions
  const { user, stats, logout, loading } = useAuth();
  // State to control the visibility of the donation modal
  const [showDonationModal, setShowDonationModal] = useState(false);

  return (
    <>
      <DonationModal show={showDonationModal} onClose={() => setShowDonationModal(false)} />
      <header className={styles.header}>
      <div className={styles.container}>
        <Link href="/" legacyBehavior>
          <a className={styles.logo}>FreeEdu</a>
        </Link>
        <div className={styles.searchContainer}>
          <SearchBar />
        </div>
        <nav>
          <ul className={styles.navList}>
            <li>
              <Link href="/" legacyBehavior>
                <a>Courses</a>
              </Link>
            </li>
            <li>
              <Link href="/leaderboard" legacyBehavior>
                <a>Leaderboard</a>
              </Link>
            </li>
            <li>
              <Link href="/paths" legacyBehavior>
                <a>Learning Paths</a>
              </Link>
            </li>
            {loading ? null : user ? (
              <>
                <li>
                  <Link href="/profile" legacyBehavior>
                    <a className={styles.userName}>
                      Hello, {user.name} (⭐ {stats.score || 0})
                    </a>
                  </Link>
                </li>
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
            <li>
              <button onClick={() => setShowDonationModal(true)} className={styles.supportButton}>Support Us</button>
            </li>
          </ul>
        </nav>
      </div>
    </header>
    </>
  );
};

export default Header;
