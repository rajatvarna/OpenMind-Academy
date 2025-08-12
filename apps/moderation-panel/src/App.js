import React, { useState } from 'react';
import './App.css';
import ModerationQueue from './components/ModerationQueue';
import SubmissionDetail from './components/SubmissionDetail';

function App() {
  // This state will hold the submission currently being viewed by the moderator.
  const [selectedSubmission, setSelectedSubmission] = useState(null);

  // This function will be passed down to the queue to update the selected submission.
  const handleSelectSubmission = (submission) => {
    setSelectedSubmission(submission);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Moderation Panel</h1>
      </header>
      <main className="App-main">
        <div className="queue-panel">
          <ModerationQueue onSelectSubmission={handleSelectSubmission} />
        </div>
        <div className="detail-panel">
          {selectedSubmission ? (
            <SubmissionDetail submission={selectedSubmission} />
          ) : (
            <p>Select a submission from the queue to review it.</p>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;
