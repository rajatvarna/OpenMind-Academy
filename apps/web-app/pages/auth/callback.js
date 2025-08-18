import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { useAuth } from '../../context/AuthContext';

export default function AuthCallbackPage() {
  const router = useRouter();
  const { login } = useAuth();

  useEffect(() => {
    if (router.isReady) {
      const { token } = router.query;
      if (token) {
        login(token);
        // Redirect to the profile page after successful login
        router.push('/profile');
      } else {
        // Handle error or no token case
        router.push('/login?error=Authentication failed');
      }
    }
  }, [router.isReady, router.query, login, router]);

  return (
    <div className="container">
      <main>
        <p>Logging you in...</p>
      </main>
    </div>
  );
}
