import React from 'react';
import Link from 'next/link';
import { useAuth } from '../context/AuthContext';
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
      <header className="w-full py-4 bg-white border-b border-gray-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center">
          <Link href="/" legacyBehavior>
            <a className="text-2xl font-bold text-gray-900">FreeEdu</a>
          </Link>
          <div className="flex-1 flex justify-center px-2 lg:ml-6 lg:justify-end">
            <SearchBar />
          </div>
          <nav className="hidden md:flex items-center space-x-8">
            <Link href="/" legacyBehavior><a className="text-gray-600 hover:text-blue-600">Courses</a></Link>
            <Link href="/leaderboard" legacyBehavior><a className="text-gray-600 hover:text-blue-600">Leaderboard</a></Link>
            <Link href="/paths" legacyBehavior><a className="text-gray-600 hover:text-blue-600">Learning Paths</a></Link>
            {loading ? null : user ? (
              <>
                <Link href="/profile" legacyBehavior>
                  <a className="text-gray-800 font-semibold">
                    Hello, {user.name} (⭐ {stats.score || 0})
                  </a>
                </Link>
                <Link href="/settings" legacyBehavior><a className="text-gray-600 hover:text-blue-600">Settings</a></Link>
                <button onClick={logout} className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
                  Logout
                </button>
              </>
            ) : (
              <Link href="/login" legacyBehavior>
                <a className="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700">
                  Login
                </a>
              </Link>
            )}
            <button
              onClick={() => setShowDonationModal(true)}
              className="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-yellow-900 bg-yellow-400 hover:bg-yellow-500"
            >
              Support Us
            </button>
          </nav>
        </div>
      </header>
    </>
  );
};

export default Header;
