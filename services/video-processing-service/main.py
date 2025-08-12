import pika
import time
import os
import json
from process_video import handle_video_request

# --- RabbitMQ Configuration ---
RABBITMQ_URL = os.environ.get('RABBITMQ_URL', 'amqp://guest:guest@rabbitmq:5672')
QUEUE_NAME = 'video_generation_queue'

def main():
    """
    Main function to start the RabbitMQ consumer.
    Includes connection retry logic.
    """
    while True:
        try:
            connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
            channel = connection.channel()

            # Declare the queue. 'durable=True' ensures it survives broker restarts.
            channel.queue_declare(queue=QUEUE_NAME, durable=True)
            print(' [*] Waiting for messages. To exit press CTRL+C')

            # This tells RabbitMQ not to give more than one message to this worker at a time.
            # This way, the worker won't be overwhelmed.
            channel.basic_qos(prefetch_count=1)

            # Set up the consumer
            channel.basic_consume(queue=QUEUE_NAME, on_message_callback=callback)

            # Start consuming
            channel.start_consuming()

        except pika.exceptions.AMQPConnectionError as e:
            print(f"Connection to RabbitMQ failed: {e}. Retrying in 5 seconds...")
            time.sleep(5)
        except Exception as e:
            print(f"An unexpected error occurred: {e}. Restarting...")
            time.sleep(5)

def callback(ch, method, properties, body):
    """
    Callback function executed when a message is received.
    """
    print(f" [x] Received message with delivery tag {method.delivery_tag}")

    try:
        # The body is a byte string, so we decode and parse it as JSON.
        message_body = json.loads(body.decode('utf-8'))
        print(f" [x] Message body: {message_body}")

        # Process the video request
        handle_video_request(message_body)

        # Acknowledge the message, telling RabbitMQ it has been successfully processed.
        print(f" [x] Done processing. Acknowledging message {method.delivery_tag}.")
        ch.basic_ack(delivery_tag=method.delivery_tag)

    except json.JSONDecodeError as e:
        print(f"Error decoding message body: {e}")
        # Reject the message as it's malformed and cannot be processed.
        # 'requeue=False' sends it to a dead-letter queue if one is configured.
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    except Exception as e:
        print(f"An error occurred while processing the message: {e}")
        # Negative acknowledgment. You might want to requeue depending on the error.
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)
