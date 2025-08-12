import '../styles/globals.css';
import Layout from '../components/Layout';

// This App component is the top-level component which will be common across all different pages.
// You can use this App component to keep state when navigating between pages, for example.
function MyApp({ Component, pageProps }) {
  return (
    <Layout>
      <Component {...pageProps} />
    </Layout>
  );
}

export default MyApp;
