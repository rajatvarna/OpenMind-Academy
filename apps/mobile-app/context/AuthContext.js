import React, { createContext, useState, useContext } from 'react';

const AuthContext = createContext();

// This is a simplified, in-memory auth provider for demonstration.
// A real React Native app would use @react-native-async-storage/async-storage
// for persistence and would handle fetching data from API endpoints.
export const AuthProvider = ({ children }) => {
  // We'll use mock data since we can't make API calls here easily.
  const [user, setUser] = useState({ id: 1, name: 'Mobile User' });
  const [stats, setStats] = useState({ score: 120 });
  const [completedLessons, setCompletedLessons] = useState(new Set([1, 3])); // Mock completed lessons

  const value = {
    user,
    stats,
    completedLessons,
    // In a real app, you'd have login/logout functions here
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook to use the auth context
export const useAuth = () => useContext(AuthContext);
