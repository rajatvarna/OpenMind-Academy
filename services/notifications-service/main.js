const amqp = require('amqplib');
const { handleEvent } = require('./handler');

const RABBITMQ_URL = process.env.RABBITMQ_URL || 'amqp://guest:guest@rabbitmq:5672';
const QUEUE_NAME = 'notifications_events';

async function main() {
  console.log('Starting Notifications Service...');

  let connection;
  try {
    connection = await amqp.connect(RABBITMQ_URL);
    const channel = await connection.createChannel();

    await channel.assertQueue(QUEUE_NAME, { durable: true });
    channel.prefetch(1); // Process one message at a time

    console.log(`[*] Waiting for messages in ${QUEUE_NAME}. To exit press CTRL+C`);

    channel.consume(QUEUE_NAME, async (msg) => {
      if (msg !== null) {
        try {
          const messageContent = msg.content.toString();
          console.log(`[x] Received: ${messageContent}`);

          const { eventType, payload } = JSON.parse(messageContent);

          await handleEvent(eventType, payload);

          // Acknowledge the message
          channel.ack(msg);
          console.log('[x] Done. Message acknowledged.');

        } catch (error) {
          console.error('Error processing message:', error);
          // Reject the message without requeueing to avoid poison pills
          channel.nack(msg, false, false);
        }
      }
    });

  } catch (error) {
    console.error('Failed to start service:', error);
    // In a real app, you'd have a more robust retry mechanism or process exit.
    setTimeout(main, 5000); // Retry connection after 5 seconds
  }
}

main();
