import { useEffect } from 'react';
import '../styles/globals.css';
import Layout from '../components/Layout';
import { AuthProvider, useAuth } from '../context/AuthContext';
import { loadStripe } from '@stripe/stripe-js';
import { Elements } from '@stripe/react-stripe-js';

// In a real application, load this from an environment variable.
const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY);

const ThemeManager = ({ children }) => {
  const { user } = useAuth();

  useEffect(() => {
    if (user && user.preferences && user.preferences.theme) {
      document.body.className = user.preferences.theme;
    }
  }, [user]);

  return children;
};

// This App component is the top-level component which will be common across all different pages.
// You can use this App component to keep state when navigating between pages, for example.
function MyApp({ Component, pageProps }) {
  return (
    <AuthProvider>
      <ThemeManager>
        <Elements stripe={stripePromise}>
          <Layout>
            <Component {...pageProps} />
          </Layout>
        </Elements>
      </ThemeManager>
    </AuthProvider>
  );
}

export default MyApp;
