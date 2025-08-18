# Proposed Feature Enhancements

Based on a review of the existing codebase, here are five proposed feature enhancements that would significantly improve the platform's functionality, security, and user experience.

---

### 1. User Profile Picture Management

**Description:**
Allow users to personalize their profiles by uploading a profile picture. This feature would involve creating new API endpoints in the `user-service` to handle image uploads, updates, and deletions. The images themselves would be stored in a cloud storage bucket (e.g., GCS or S3), and the `user-service` would store a reference to the image URL in the `users` table.

**Components Affected:**
- `user-service`: New endpoints for upload, delete. New field in `users` model.
- `web-app` / `mobile-app`: UI components for uploading and displaying the profile picture.
- `infrastructure`: Potentially new or updated cloud storage bucket policies.

---

### 2. Two-Factor Authentication (2FA)

**Description:**
Greatly enhance account security by implementing Time-based One-Time Password (TOTP) as a second factor of authentication. When a user enables 2FA, they would use an authenticator app (like Google Authenticator) to generate a code required for login. This would also include the generation of single-use recovery codes.

**Components Affected:**
- `user-service`: Endpoints to enable/disable 2FA, validate TOTP codes, and manage recovery codes.
- `web-app` / `mobile-app`: UI flows for setting up 2FA and for prompting the user for a code during login.
- `api-gateway`: The login flow would need to be updated to handle the 2FA challenge.

---

### 3. User Account Deactivation and Deletion

**Description:**
Provide users with more control over their data by allowing them to either temporarily deactivate their account (a soft delete) or permanently delete it (a hard delete). Deactivation would make the profile inaccessible but preserve the data, while deletion would remove all user data from the system, possibly after a grace period.

**Components Affected:**
- `user-service`: New endpoints to handle deactivation and deletion requests. Logic to anonymize or erase user data.
- `web-app` / `mobile-app`: UI elements in the user settings page to initiate these actions.
- All other services: A mechanism (e.g., event-based) would be needed to propagate the deletion to other services that hold user data.

---

### 4. Detailed User Activity Log

**Description:**
Create a new API endpoint that returns a log of a user's recent activities, such as courses enrolled in, lessons completed, comments posted, and badges earned. This would provide users with a clear history of their engagement on the platform and could be used to build a more detailed profile page.

**Components Affected:**
- `user-service`: Would need a new table to store activity events and a new endpoint to retrieve them.
- All other services: Services would need to publish events (e.g., to RabbitMQ) whenever a user performs a relevant action. The `user-service` would consume these events and populate the activity log.
- `web-app` / `mobile-app`: A new UI section on the profile page to display the activity feed.

---

### 5. Social Login Integration (OAuth 2.0)

**Description:**
Streamline the registration and login process by allowing users to sign in with their existing accounts from third-party providers like Google, GitHub, or Facebook. This would lower the barrier to entry for new users and simplify the login experience for existing ones.

**Components Affected:**
- `user-service`: New endpoints to handle the OAuth 2.0 callback from providers, create a new user or link to an existing one, and issue a JWT.
- `web-app` / `mobile-app`: UI buttons and logic to initiate the OAuth 2.0 flow.
- `api-gateway`: The gateway would need to be configured to handle the new authentication routes.
