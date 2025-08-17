import { createContext, useState, useEffect, useContext } from 'react';
import { useRouter } from 'next/router';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [stats, setStats] = useState({ score: 0 });
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const checkUser = async () => {
      try {
        const userRes = await fetch('/api/me');
        if (userRes.ok) {
          const { user } = await userRes.json();
          setUser(user);
          // If user exists, fetch their stats
          const statsRes = await fetch(`/api/stats/${user.id}`);
          if (statsRes.ok) {
            const statsData = await statsRes.json();
            setStats(statsData);
          }
        } else {
          setUser(null);
          setStats({ score: 0 });
        }
      } catch (error) {
        setUser(null);
        setStats({ score: 0 });
      } finally {
        setLoading(false);
      }
    };
    checkUser();
  }, []);

  const logout = async () => {
    // Create a logout API route to clear the cookie
    try {
      await fetch('/api/logout');
      setUser(null);
      router.push('/'); // Redirect to homepage after logout
    } catch (error) {
      console.error('Failed to logout', error);
    }
  };

  return (
    <AuthContext.Provider value={{ user, stats, loading, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook to use the auth context
export const useAuth = () => useContext(AuthContext);
