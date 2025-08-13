import { useState, useEffect } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { useAuth } from '../context/AuthContext';
import styles from '../styles/Profile.module.css';

export default function ProfilePage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [profileData, setProfileData] = useState(null);
  const [isLoadingProfile, setIsLoadingProfile] = useState(true);

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login');
    }
    if (user) {
      const fetchProfile = async () => {
        try {
          setIsLoadingProfile(true);
          const res = await fetch(`/api/profile/${user.id}`);
          if (res.ok) {
            const data = await res.json();
            setProfileData(data);
          }
        } catch (error) {
          console.error("Failed to fetch profile data", error);
        } finally {
          setIsLoadingProfile(false);
        }
      };
      fetchProfile();
    }
  }, [user, loading, router]);

  if (loading || isLoadingProfile) {
    return <div className="container">Loading profile...</div>;
  }

  if (!profileData) {
    return <div className="container">Could not load profile data.</div>;
  }

  return (
    <div className="container">
      <Head>
        <title>Your Profile</title>
      </Head>
      <main>
        <h1>{user.name}'s Profile</h1>
        <div className={styles.profileGrid}>
          <div className={styles.statCard}>
            <h3>Total Score</h3>
            <p>{profileData.gamification_stats?.score || 0}</p>
          </div>
          <div className={styles.statCard}>
            <h3>Courses Created</h3>
            <p>{profileData.created_courses?.length || 0}</p>
          </div>
        </div>

        <div className={styles.coursesSection}>
          <h2>My Courses</h2>
          {/* Here we would map over profileData.created_courses and display them */}
          <p>A list of your created courses would appear here.</p>
        </div>
      </main>
    </div>
  );
}
