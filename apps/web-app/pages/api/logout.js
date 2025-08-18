export default async function handler(req, res) {
  // The cookie is now managed on the client-side.
  // This endpoint can be kept for future use (e.g., server-side session invalidation)
  // but for now, it does nothing.
  res.status(200).json({ success: true, message: 'Logged out successfully' });
}
