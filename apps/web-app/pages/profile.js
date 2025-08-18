import { useState, useEffect } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { useAuth } from '../context/AuthContext';

export default function ProfilePage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [profileData, setProfileData] = useState(null);
  const [quizAttempts, setQuizAttempts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login');
      return;
    }
    if (user) {
      const fetchData = async () => {
        setIsLoading(true);
        try {
          const [profileRes, attemptsRes] = await Promise.all([
            fetch(`/api/profile/${user.id}`),
            fetch(`/api/profile/quiz-attempts?userId=${user.id}`),
          ]);

          if (profileRes.ok) {
            const data = await profileRes.json();
            setProfileData(data);
          }
          if (attemptsRes.ok) {
            const data = await attemptsRes.json();
            setQuizAttempts(data);
          }
        } catch (error) {
          console.error("Failed to fetch profile data", error);
        } finally {
          setIsLoading(false);
        }
      };
      fetchData();
    }
  }, [user, loading, router]);

  if (isLoading) {
    return <div className="container mx-auto px-4">Loading profile...</div>;
  }

  if (!profileData) {
    return <div className="container mx-auto px-4">Could not load profile data.</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Head>
        <title>Your Profile</title>
      </Head>
      <main>
        <div className="flex items-center space-x-6 mb-8">
          <img
            src={user.profile_picture_url || 'https://via.placeholder.com/128'}
            alt="Profile Picture"
            className="w-32 h-32 rounded-full"
          />
          <h1 className="text-4xl font-bold">{user.name}'s Profile</h1>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <div className="p-6 bg-white rounded-lg shadow-md text-center">
            <h3 className="text-xl font-semibold text-gray-700">Total Score</h3>
            <p className="text-3xl font-bold mt-2">{profileData.gamification_stats?.score || 0}</p>
          </div>
          <div className="p-6 bg-white rounded-lg shadow-md text-center">
            <h3 className="text-xl font-semibold text-gray-700">Courses Created</h3>
            <p className="text-3xl font-bold mt-2">{profileData.created_courses?.length || 0}</p>
          </div>
          <div className="p-6 bg-white rounded-lg shadow-md text-center">
            <h3 className="text-xl font-semibold text-gray-700">Quizzes Taken</h3>
            <p className="text-3xl font-bold mt-2">{quizAttempts.length}</p>
          </div>
        </div>

        <div className="mb-12">
          <h2 className="text-3xl font-bold mb-6">Quiz History</h2>
          <div className="bg-white rounded-lg shadow-md">
            <ul className="divide-y divide-gray-200">
              {quizAttempts.length > 0 ? (
                quizAttempts.map(attempt => (
                  <li key={attempt.id} className="p-4 flex justify-between items-center">
                    <div>
                      <p className="font-semibold">Quiz ID: {attempt.quiz_id}</p>
                      <p className="text-sm text-gray-600">
                        Completed on: {new Date(attempt.created_at).toLocaleDateString()}
                      </p>
                    </div>
                    <div className="text-lg font-bold">
                      Score: {attempt.score}
                    </div>
                  </li>
                ))
              ) : (
                <li className="p-4 text-center text-gray-500">You haven't taken any quizzes yet.</li>
              )}
            </ul>
          </div>
        </div>

        <div>
          <h2 className="text-3xl font-bold mb-6">My Courses</h2>
          <div className="bg-white rounded-lg shadow-md p-4 text-center text-gray-500">
            {/* Here we would map over profileData.created_courses and display them */}
            <p>A list of your created courses would appear here.</p>
          </div>
        </div>
      </main>
    </div>
  );
}
