import { useState } from 'react';

export default function ProfilePictureUploader({ onUploadSuccess }) {
  const [file, setFile] = useState(null);
  const [error, setError] = useState(null);
  const [message, setMessage] = useState(null);
  const [isUploading, setIsUploading] = useState(false);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
    setError(null);
    setMessage(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file) {
      setError('Please select a file to upload.');
      return;
    }

    setIsUploading(true);
    setError(null);
    setMessage(null);

    const formData = new FormData();
    formData.append('picture', file);

    try {
      const res = await fetch('/api/users/profile/picture', {
        method: 'POST',
        // No 'Content-Type' header needed, the browser will set it
        // correctly for multipart/form-data.
        body: formData,
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || 'Something went wrong');
      }

      setMessage(data.message);
      if (onUploadSuccess) {
        onUploadSuccess(data.url);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div>
      <h3>Upload New Profile Picture</h3>
      <form onSubmit={handleSubmit}>
        <input type="file" accept="image/png, image/jpeg" onChange={handleFileChange} />
        <button type="submit" disabled={!file || isUploading}>
          {isUploading ? 'Uploading...' : 'Upload'}
        </button>
      </form>
      {error && <p style={{ color: 'red' }}>Error: {error}</p>}
      {message && <p style={{ color: 'green' }}>{message}</p>}
    </div>
  );
}
