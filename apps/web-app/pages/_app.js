import '../styles/globals.css';
import Layout from '../components/Layout';
import { AuthProvider } from '../context/AuthContext';
import { loadStripe } from '@stripe/stripe-js';
import { Elements } from '@stripe/react-stripe-js';

// Use a placeholder publishable key.
// In a real application, load this from an environment variable.
const stripePromise = loadStripe('pk_test_51...PLACEHOLDER');

// This App component is the top-level component which will be common across all different pages.
// You can use this App component to keep state when navigating between pages, for example.
function MyApp({ Component, pageProps }) {
  return (
    <AuthProvider>
      <Elements stripe={stripePromise}>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </Elements>
    </AuthProvider>
  );
}

export default MyApp;
