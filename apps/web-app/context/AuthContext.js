import { createContext, useState, useEffect, useContext } from 'react';
import { useRouter } from 'next/router';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [stats, setStats] = useState({ score: 0 });
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const login = (token) => {
    // In a real app, you'd use a library like js-cookie
    document.cookie = `auth_token=${token}; path=/; max-age=86400`; // 1 day
    checkUser();
  };

  const checkUser = async () => {
    const token = document.cookie.split('; ').find(row => row.startsWith('auth_token='))?.split('=')[1];
    if (!token) {
      setUser(null);
      setLoading(false);
      return;
    }

    try {
      const userRes = await fetch('/api/me'); // This API route will now need to read the cookie
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

  useEffect(() => {
    checkUser();
  }, []);

  const logout = () => {
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    setUser(null);
    setStats({ score: 0 });
    router.push('/'); // Redirect to homepage after logout
  };

  const refetchUser = async () => {
    await checkUser();
  };

  return (
    <AuthContext.Provider value={{ user, stats, loading, login, logout, refetchUser }}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook to use the auth context
export const useAuth = () => useContext(AuthContext);
