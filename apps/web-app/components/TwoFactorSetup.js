import { useState } from 'react';
import { useAuth } from '../context/AuthContext';

// In a real app, you would install this with: npm install qrcode.react
// import QRCode from 'qrcode.react';

export default function TwoFactorSetup() {
  const { user, refetchUser } = useAuth();
  const [setupData, setSetupData] = useState(null);
  const [token, setToken] = useState('');
  const [error, setError] = useState(null);
  const [message, setMessage] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleEnable = async () => {
    setIsLoading(true);
    setError(null);
    setMessage(null);
    try {
      const res = await fetch('/api/users/2fa/enable', { method: 'POST' });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to start 2FA setup.');
      setSetupData(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerify = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    setMessage(null);
    try {
      const res = await fetch('/api/users/2fa/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to verify token.');
      setMessage(data.message);
      setSetupData(null); // Clear setup data on success
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDisable = async () => {
    setIsLoading(true);
    setError(null);
    setMessage(null);
    try {
      const res = await fetch('/api/users/2fa/disable', { method: 'POST' });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to disable 2FA.');
      setMessage(data.message);
      await refetchUser();
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (user && user.two_factor_enabled) {
    return (
      <div>
        <p className="text-green-600">Two-Factor Authentication is currently enabled.</p>
        <button onClick={handleDisable} disabled={isLoading} className="mt-2 text-red-600">
          {isLoading ? 'Disabling...' : 'Disable 2FA'}
        </button>
        {error && <p className="text-red-600 mt-2">Error: {error}</p>}
        {message && <p className="text-green-600 mt-2">{message}</p>}
      </div>
    );
  }

  if (!setupData) {
    return (
      <button onClick={handleEnable} disabled={isLoading}>
        {isLoading ? 'Loading...' : 'Enable Two-Factor Authentication'}
      </button>
    );
  }

  return (
    <div>
      <h4 className="font-bold">Step 1: Scan this QR Code</h4>
      <p>Scan the image below with your authenticator app (e.g., Google Authenticator).</p>
      <img
        src={`https://chart.googleapis.com/chart?chs=200x200&cht=qr&chl=${encodeURIComponent(setupData.otpauth_url)}`}
        alt="QR Code"
      />

      <h4 className="font-bold mt-4">Step 2: Save Your Recovery Codes</h4>
      <p>Store these codes in a safe place. They can be used to access your account if you lose your device.</p>
      <div className="p-4 bg-gray-100 my-2">
        <ul className="list-disc list-inside">
          {setupData.recovery_codes.map(code => <li key={code}><code>{code}</code></li>)}
        </ul>
      </div>

      <h4 className="font-bold mt-4">Step 3: Verify Your Device</h4>
      <p>Enter the 6-digit code from your authenticator app to complete the setup.</p>
      <form onSubmit={handleVerify} className="flex items-center space-x-2 mt-2">
        <input
          type="text"
          value={token}
          onChange={(e) => setToken(e.target.value)}
          placeholder="6-digit code"
          className="border p-2"
          maxLength="6"
        />
        <button type="submit" disabled={isLoading || !token}>
          {isLoading ? 'Verifying...' : 'Verify & Activate'}
        </button>
      </form>
      {error && <p className="text-red-600 mt-2">Error: {error}</p>}
      {message && <p className="text-green-600 mt-2">{message}</p>}
    </div>
  );
}
