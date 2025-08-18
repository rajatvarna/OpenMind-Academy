import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/Settings.module.css';
import ProfilePictureUploader from '../components/ProfilePictureUploader';
import TwoFactorSetup from '../components/TwoFactorSetup';
import Modal from '../components/Modal';

export default function SettingsPage() {
  const { user, refetchUser, logout } = useAuth();
  const [preferences, setPreferences] = useState({ theme: 'light' });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (user) {
      // Fetch initial preferences
      const fetchPreferences = async () => {
        try {
          const res = await fetch('/api/profile/preferences');
          if (res.ok) {
            const data = await res.json();
            setPreferences(data);
          } else {
            setError('Could not fetch preferences.');
          }
        } catch (err) {
          setError('An error occurred while fetching preferences.');
        } finally {
          setLoading(false);
        }
      };
      fetchPreferences();
    }
  }, [user]);

  const handleThemeChange = async (e) => {
    const newTheme = e.target.value;
    setPreferences({ ...preferences, theme: newTheme });

    try {
      await fetch('/api/profile/preferences', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ theme: newTheme }),
      });
      await refetchUser();
    } catch (err) {
      setError('Failed to save theme preference.');
    }
  };

  const handleUploadSuccess = () => {
    // Refetch user data to get the new profile picture URL
    refetchUser();
  };

  const [showDeactivateModal, setShowDeactivateModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const handleDeactivate = async () => {
    setError(null);
    try {
      const res = await fetch('/api/users/profile', { method: 'DELETE' });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || 'Failed to deactivate account.');
      }
      // Logout will handle redirecting the user
      await logout();
      setShowDeactivateModal(false);
    } catch (err) {
      setError(err.message);
      setShowDeactivateModal(false);
    }
  };

  const handleDelete = async () => {
    setError(null);
    try {
      const res = await fetch('/api/users/account', { method: 'DELETE' });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || 'Failed to delete account.');
      }
      // Logout will handle redirecting the user
      await logout();
      setShowDeleteModal(false);
    } catch (err) {
      setError(err.message);
      setShowDeleteModal(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div className={styles.error}>{error}</div>;
  }

  return (
    <div className="container">
      <h1 className={styles.title}>Settings</h1>
      <div className={styles.settingsForm}>
        <h2>Theme</h2>
        <p>Choose your preferred theme for the application.</p>
        <select value={preferences.theme} onChange={handleThemeChange}>
          <option value="light">Light</option>
          <option value="dark">Dark</option>
        </select>
      </div>

      <div className={styles.settingsForm}>
        <h2>Profile Picture</h2>
        <p>Upload a new profile picture.</p>
        <ProfilePictureUploader onUploadSuccess={handleUploadSuccess} />
      </div>

      <div className={styles.settingsForm}>
        <h2>Two-Factor Authentication (2FA)</h2>
        <p>Add an extra layer of security to your account.</p>
        <TwoFactorSetup />
      </div>

      <div className={`${styles.settingsForm} ${styles.dangerZone}`}>
        <h2>Danger Zone</h2>
        <div className={styles.dangerAction}>
          <p>Deactivate your account. Your data will be kept but your profile will not be public.</p>
          <button
            onClick={() => setShowDeactivateModal(true)}
            className={styles.button}
          >
            Deactivate Account
          </button>
        </div>
        <div className={styles.dangerAction}>
          <p>Permanently delete your account and all of your data. This action is irreversible.</p>
          <button
            onClick={() => setShowDeleteModal(true)}
            className={styles.dangerButton}
          >
            Delete Account
          </button>
        </div>
      </div>

      <Modal
        show={showDeactivateModal}
        onClose={() => setShowDeactivateModal(false)}
        title="Deactivate Account"
      >
        <p>Are you sure you want to deactivate your account? This action cannot be undone.</p>
        <div className={styles.modalActions}>
          <button onClick={() => setShowDeactivateModal(false)} className={styles.button}>
            Cancel
          </button>
          <button onClick={handleDeactivate} className={styles.dangerButton}>
            Deactivate
          </button>
        </div>
      </Modal>

      <Modal
        show={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        title="Permanently Delete Account"
      >
        <p>Are you absolutely sure? This will permanently delete your account, profile, and all associated data. This action cannot be undone.</p>
        <div className={styles.modalActions}>
          <button onClick={() => setShowDeleteModal(false)} className={styles.button}>
            Cancel
          </button>
          <button onClick={handleDelete} className={styles.dangerButton}>
            Yes, Delete My Account
          </button>
        </div>
      </Modal>
    </div>
  );
}
