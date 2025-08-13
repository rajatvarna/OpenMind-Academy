import { useState } from 'react';
import styles from '../styles/ChatInterface.module.css';

export default function ChatInterface() {
  const [messages, setMessages] = useState([
    { from: 'bot', text: 'Hello! Ask me anything about this lesson.' }
  ]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSend = async (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    const userMessage = { from: 'user', text: input };
    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    try {
      // Send the question to our backend API route
      const res = await fetch('/api/qna', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ question: input }),
      });

      if (res.ok) {
        const data = await res.json();
        setMessages(prev => [...prev, { from: 'bot', text: data.answer }]);
      } else {
        setMessages(prev => [...prev, { from: 'bot', text: 'Sorry, I had trouble finding an answer.' }]);
      }
    } catch (error) {
      setMessages(prev => [...prev, { from: 'bot', text: 'Sorry, there was an error connecting to the service.' }]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={styles.chatContainer}>
      <div className={styles.messageList}>
        {messages.map((msg, index) => (
          <div key={index} className={`${styles.message} ${styles[msg.from]}`}>
            {msg.text}
          </div>
        ))}
        {isLoading && <div className={`${styles.message} ${styles.bot}`}>Thinking...</div>}
      </div>
      <form onSubmit={handleSend} className={styles.inputForm}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Ask a question..."
          className={styles.input}
          disabled={isLoading}
        />
        <button type="submit" className={styles.sendButton} disabled={isLoading}>
          Send
        </button>
      </form>
    </div>
  );
}
