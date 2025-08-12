const sgMail = require('@sendgrid/mail');
const admin = require('firebase-admin');

// --- SendGrid Setup ---
// In a real app, this API key would be loaded securely from env vars or a secret manager.
const SENDGRID_API_KEY = process.env.SENDGRID_API_KEY || 'YOUR_SENDGRID_API_KEY_PLACEHOLDER';
sgMail.setApiKey(SENDGRID_API_KEY);

// --- Firebase Admin SDK Setup ---
// In a real app, you would initialize this with service account credentials.
// const serviceAccount = require('./path/to/your/serviceAccountKey.json');
// admin.initializeApp({
//   credential: admin.credential.cert(serviceAccount)
// });
let firebaseInitialized = false;
if (process.env.FIREBASE_SERVICE_ACCOUNT) {
    try {
        const serviceAccount = JSON.parse(process.env.FIREBASE_SERVICE_ACCOUNT);
        admin.initializeApp({
            credential: admin.credential.cert(serviceAccount)
        });
        firebaseInitialized = true;
        console.log("Firebase Admin SDK initialized successfully.");
    } catch (e) {
        console.error("Failed to initialize Firebase Admin SDK:", e.message);
    }
} else {
    console.log("Firebase credentials not found. Push notifications will be disabled.");
}


/**
 * Sends an email using SendGrid.
 * @param {string} to - The recipient's email address.
 * @param {string} subject - The subject of the email.
 * @param {string} html - The HTML body of the email.
 */
async function sendEmail({ to, subject, html }) {
    const msg = {
        to,
        from: 'no-reply@freeedu.com', // This should be a verified sender in your SendGrid account
        subject,
        html,
    };

    console.log('--- SIMULATING EMAIL ---');
    console.log(`To: ${to}`);
    console.log(`Subject: ${subject}`);
    console.log(`Body: ${html}`);
    console.log('------------------------');

    // In a real app with a valid API key, you would uncomment the following lines:
    // try {
    //     await sgMail.send(msg);
    //     console.log(`Email sent to ${to}`);
    // } catch (error) {
    //     console.error('Error sending email:', error);
    // }
    return Promise.resolve();
}

/**
 * Sends a push notification using Firebase Cloud Messaging.
 * @param {string} deviceToken - The FCM token for the target device.
 * @param {string} title - The title of the notification.
 * @param {string} body - The body of the notification.
 */
async function sendPushNotification({ deviceToken, title, body }) {
    const message = {
        notification: {
            title,
            body,
        },
        token: deviceToken,
    };

    console.log('--- SIMULATING PUSH NOTIFICATION ---');
    console.log(`To Device: ${deviceToken}`);
    console.log(`Title: ${title}`);
    console.log(`Body: ${body}`);
    console.log('------------------------------------');

    // In a real app with valid credentials, you would uncomment the following:
    // if (!firebaseInitialized) {
    //     console.error("Cannot send push notification, Firebase not initialized.");
    //     return;
    // }
    // try {
    //     const response = await admin.messaging().send(message);
    //     console.log('Successfully sent message:', response);
    // } catch (error) {
    //     console.error('Error sending message:', error);
    // }
    return Promise.resolve();
}

module.exports = {
    sendEmail,
    sendPushNotification,
};
