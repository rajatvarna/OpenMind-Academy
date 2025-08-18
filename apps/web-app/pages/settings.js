import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/Settings.module.css';
import ProfilePictureUploader from '../components/ProfilePictureUploader';
import TwoFactorSetup from '../components/TwoFactorSetup';

export default function SettingsPage() {
  const { user, refetchUser } = useAuth();
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
    </div>
  );
}
