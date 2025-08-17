import { useState } from 'react';

export default function Quiz({ quizData }) {
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);
  const [showScore, setShowScore] = useState(false);
  const [selectedAnswer, setSelectedAnswer] = useState(null);
  const [isCorrect, setIsCorrect] = useState(null);

  const currentQuestion = quizData.questions[currentQuestionIndex];

  const handleAnswerClick = (answerIndex) => {
    setSelectedAnswer(answerIndex);
    const correct = answerIndex === currentQuestion.correctAnswer;
    setIsCorrect(correct);
    if (correct) {
      setScore(prev => prev + 1);
    }
  };

  const handleNextQuestion = () => {
    setSelectedAnswer(null);
    setIsCorrect(null);
    const nextQuestion = currentQuestionIndex + 1;
    if (nextQuestion < quizData.questions.length) {
      setCurrentQuestionIndex(nextQuestion);
    } else {
      setShowScore(true);
    }
  };

  if (showScore) {
    return (
      <div className="p-8 bg-white rounded-lg shadow-lg max-w-2xl mx-auto text-center">
        <h2 className="text-2xl font-bold mb-4">Quiz Complete!</h2>
        <p className="text-xl">You scored {score} out of {quizData.questions.length}</p>
      </div>
    );
  }

  return (
    <div className="p-8 bg-white rounded-lg shadow-lg max-w-2xl mx-auto">
      <div className="mb-6">
        <div className="flex justify-between items-center mb-2 text-gray-600">
          <span>Question {currentQuestionIndex + 1}</span>
          <span>/{quizData.questions.length}</span>
        </div>
        <h2 className="text-2xl font-bold">{currentQuestion.question}</h2>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {currentQuestion.options.map((option, index) => {
          const isSelected = selectedAnswer === index;
          let buttonClass = "w-full p-4 border rounded-lg text-left transition-colors ";
          if (isSelected) {
            buttonClass += isCorrect ? "bg-green-200 border-green-500" : "bg-red-200 border-red-500";
          } else {
            buttonClass += "hover:bg-gray-100";
          }
          return (
            <button
              key={index}
              className={buttonClass}
              onClick={() => handleAnswerClick(index)}
              disabled={selectedAnswer !== null}
            >
              {option}
            </button>
          );
        })}
      </div>
      {selectedAnswer !== null && (
        <div className="mt-6 text-center">
          <button
            onClick={handleNextQuestion}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
