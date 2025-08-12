const express = require('express');
const amqp = require('amqplib');
const helmet = require('helmet');
const cors = require('cors');
const morgan = require('morgan');

const app = express();
const PORT = process.env.PORT || 3002;

// Middleware
app.use(helmet());
app.use(cors());
app.use(morgan('combined')); // Request logging
app.use(express.json());

// --- RabbitMQ Connection and Publishing ---
const RABBITMQ_URL = process.env.RABBITMQ_URL || 'amqp://guest:guest@rabbitmq:5672';
const QUEUE_NAME = 'video_generation_queue';

let channel = null;

async function connectToRabbitMQ() {
    try {
        // In a real app, you'd have more robust connection logic with retries.
        const connection = await amqp.connect(RABBITMQ_URL);
        channel = await connection.createChannel();
        await channel.assertQueue(QUEUE_NAME, { durable: true }); // Ensure queue exists and survives restarts
        console.log('Successfully connected to RabbitMQ and asserted queue.');
    } catch (error) {
        console.error('Failed to connect to RabbitMQ:', error);
        // In a real app, you might want to exit or implement a retry mechanism.
        // For this sandbox, we'll just log the error.
        channel = null; // Ensure channel is null if connection fails
    }
}

async function publishToQueue(message) {
    if (!channel) {
        console.error('Cannot publish message: RabbitMQ channel is not available.');
        // In a real app, you might return an error to the user or try to reconnect.
        throw new Error('Message queue not available');
    }

    // The message is sent as a Buffer. 'persistent: true' makes the message survive a broker restart.
    channel.sendToQueue(QUEUE_NAME, Buffer.from(JSON.stringify(message)), { persistent: true });
    console.log(`Message sent to queue '${QUEUE_NAME}':`, message);
}


// --- API Endpoint ---
app.post('/api/v1/submit', async (c, res) => {
    const { textContent, lessonId } = c.body;

    if (!textContent || !lessonId) {
        return res.status(400).json({ error: 'textContent and lessonId are required.' });
    }

    // This is where you might first save the submission to your own database
    // with a 'PENDING_VIDEO' status. For this service, we'll assume the lesson
    // was already created in the Content Service and we're just adding a video to it.

    const message = {
        lessonId,
        textContent,
        submittedAt: new Date().toISOString()
    };

    try {
        await publishToQueue(message);
        res.status(202).json({ status: 'accepted', message: 'Content submission accepted for video generation.' });
    } catch (error) {
        res.status(500).json({ error: 'Failed to queue content for processing.' });
    }
});

// Health check endpoint
app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

// --- Server Startup ---
app.listen(PORT, () => {
    console.log(`UGC Submission Service listening on port ${PORT}`);
    connectToRabbitMQ();
});
