import { useState } from 'react';
import styles from '../styles/Quiz.module.css';

export default function Quiz({ quizData }) {
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);
  const [showScore, setShowScore] = useState(false);

  const handleAnswerClick = (isCorrect) => {
    if (isCorrect) {
      setScore(prev => prev + 1);
    }

    const nextQuestion = currentQuestionIndex + 1;
    if (nextQuestion < quizData.questions.length) {
      setCurrentQuestionIndex(nextQuestion);
    } else {
      setShowScore(true);
    }
  };

  if (showScore) {
    return (
      <div className={styles.quizContainer}>
        <div className={styles.scoreSection}>
          You scored {score} out of {quizData.questions.length}
        </div>
      </div>
    );
  }

  const currentQuestion = quizData.questions[currentQuestionIndex];

  return (
    <div className={styles.quizContainer}>
      <div className={styles.questionSection}>
        <div className={styles.questionCount}>
          <span>Question {currentQuestionIndex + 1}</span>/{quizData.questions.length}
        </div>
        <div className={styles.questionText}>{currentQuestion.question}</div>
      </div>
      <div className={styles.answerSection}>
        {currentQuestion.options.map((option, index) => (
          <button
            key={index}
            onClick={() => handleAnswerClick(index === currentQuestion.correctAnswer)}
          >
            {option}
          </button>
        ))}
      </div>
    </div>
  );
}
