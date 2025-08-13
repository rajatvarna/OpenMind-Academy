import React, { useState, useEffect } from 'react';
import './ReportsQueue.css';

export default function ReportsQueue() {
  const [reports, setReports] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        setIsLoading(true);
        // In a real app, this would hit the API gateway
        // which would proxy to the UGC service.
        // const res = await fetch('/api/ugc/reports');
        // const data = await res.json();
        // setReports(data);

        // Placeholder data for now
        const placeholderReports = [
            { reportId: 1, contentId: 101, reason: 'Inaccurate information.', userId: 45 },
            { reportId: 2, contentId: 102, reason: 'Copyright violation.', userId: 88 },
            { reportId: 3, contentId: 105, reason: 'Spam or advertising.', userId: 23 },
        ];
        setReports(placeholderReports);

      } catch (error) {
        console.error('Failed to fetch reports', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchReports();
  }, []);

  if (isLoading) return <p>Loading reports...</p>;

  return (
    <div className="reports-queue">
      <h2>Pending Reports</h2>
      <table>
        <thead>
          <tr>
            <th>Report ID</th>
            <th>Content ID</th>
            <th>Reason</th>
            <th>Reported By User</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {reports.map(report => (
            <tr key={report.reportId}>
              <td>{report.reportId}</td>
              <td>{report.contentId}</td>
              <td>{report.reason}</td>
              <td>{report.userId}</td>
              <td>
                <button className="action-btn">Dismiss</button>
                <button className="action-btn-primary">Take Action</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
