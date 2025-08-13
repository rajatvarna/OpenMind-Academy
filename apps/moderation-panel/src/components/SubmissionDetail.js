import React from 'react';
import './SubmissionDetail.css';

function SubmissionDetail({ submission }) {
  if (!submission) {
    return null; // Don't render anything if no submission is selected
  }

  const handleApprove = () => {
    console.log(`Approving submission ${submission.id}`);
    // In a real app, you would make an API call here:
    // fetch(`/api/ugc/submissions/${submission.id}/approve`, { method: 'POST' });
    alert(`Submission ${submission.id} approved!`);
  };

  const handleReject = () => {
    console.log(`Rejecting submission ${submission.id}`);
    // In a real app, you would make an API call here:
    // fetch(`/api/ugc/submissions/${submission.id}/reject`, { method: 'POST' });
    alert(`Submission ${submission.id} rejected!`);
  };

  const handleDelete = () => {
    if (window.confirm(`Are you sure you want to permanently delete submission ${submission.id}? This cannot be undone.`)) {
      console.log(`Deleting submission ${submission.id}`);
      // In a real app, you would make an API call here:
      // fetch(`/api/content/courses/${submission.id}`, { method: 'DELETE' });
      alert(`Submission ${submission.id} deleted!`);
    }
  };

  return (
    <div className="submission-detail">
      <h2>{submission.title}</h2>
      <div className="author-info">
        <span>By: {submission.author}</span>
        <span>ID: {submission.id}</span>
      </div>

      <div className="content-section">
        <h3>Video Content</h3>
        {/* The video URL would come from the submission object */}
        <video controls width="100%">
          {/* In a real app, the src would be something like submission.videoUrl */}
          <source src="https://www.w3schools.com/html/mov_bbb.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
      </div>

      <div className="content-section">
        <h3>Text Content</h3>
        <p>
          This is where the text content accompanying the video would be displayed.
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus.
          Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies
          sed, dolor. Cras elementum ultrices diam. Maecenas ligula massa, varius
          a, semper congue, euismod non, mi.
        </p>
      </div>

      <div className="actions">
        <button className="approve-btn" onClick={handleApprove}>Approve</button>
        <button className="reject-btn" onClick={handleReject}>Reject</button>
      </div>
      <div className="actions destructive">
        <button className="delete-btn" onClick={handleDelete}>Delete Content Permanently</button>
      </div>
    </div>
  );
}

export default SubmissionDetail;
