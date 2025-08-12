const { sendEmail, sendPushNotification } = require('./senders');

/**
 * Main event handler function.
 * @param {string} eventType - The type of the event (e.g., 'user_registered').
 * @param {object} payload - The data associated with the event.
 */
async function handleEvent(eventType, payload) {
  console.log(`Handling event: ${eventType}`);

  switch (eventType) {
    case 'user_registered':
      return handleUserRegistered(payload);

    case 'content_approved':
      return handleContentApproved(payload);

    case 'content_rejected':
      return handleContentRejected(payload);

    default:
      console.log(`No handler for event type: ${eventType}`);
      return Promise.resolve();
  }
}

// --- Specific Event Handlers ---

/**
 * Handles the 'user_registered' event.
 * @param {object} payload - Expected to contain { email, name }.
 */
function handleUserRegistered(payload) {
  const { email, name } = payload;
  if (!email || !name) {
    console.error('Invalid payload for user_registered:', payload);
    return;
  }

  return sendEmail({
    to: email,
    subject: 'Welcome to the Free Education Platform!',
    html: `<strong>Hi ${name},</strong><p>Welcome! We're excited to have you on board.</p>`,
  });
}

/**
 * Handles the 'content_approved' event.
 * @param {object} payload - Expected to contain { authorDeviceToken, courseTitle }.
 */
function handleContentApproved(payload) {
  const { authorDeviceToken, courseTitle } = payload;
  if (!authorDeviceToken || !courseTitle) {
    console.error('Invalid payload for content_approved:', payload);
    return;
  }

  return sendPushNotification({
    deviceToken: authorDeviceToken,
    title: 'Your content is live!',
    body: `Congratulations! Your course "${courseTitle}" has been approved and is now live on the platform.`,
  });
}

/**
 * Handles the 'content_rejected' event.
 * @param {object} payload - Expected to contain { authorEmail, courseTitle, reason }.
 */
function handleContentRejected(payload) {
  const { authorEmail, courseTitle, reason } = payload;
  if (!authorEmail || !courseTitle) {
    console.error('Invalid payload for content_rejected:', payload);
    return;
  }

  return sendEmail({
    to: authorEmail,
    subject: `Update on your submission: "${courseTitle}"`,
    html: `<p>Hi there,</p><p>Thank you for your submission. Unfortunately, your course "${courseTitle}" was not approved at this time.</p><p><b>Reason:</b> ${reason || 'No reason provided.'}</p>`,
  });
}

module.exports = { handleEvent };
