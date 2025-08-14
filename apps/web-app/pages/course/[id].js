import Head from 'next/head';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import styles from '../../styles/CoursePage.module.css';
import ChatInterface from '../../components/ChatInterface';
import Reviews from '../../components/Reviews';
import ReviewForm from '../../components/ReviewForm';
import DiscussionTab from '../../components/DiscussionTab';
import ReportModal from '../../components/ReportModal';
import Modal from '../../components/Modal';
import Quiz from '../../components/Quiz';

export default function CoursePage({ course, lessons }) {
  const router = useRouter();
  const { user } = useAuth();
  const [completedLessons, setCompletedLessons] = useState(new Set());
  const [showReportModal, setShowReportModal] = useState(false);
  const [showQuizModal, setShowQuizModal] = useState(false);
  const [quizData, setQuizData] = useState(null);
  const [isLoadingQuiz, setIsLoadingQuiz] = useState(false);
  const [visibleTranscript, setVisibleTranscript] = useState(null);
  const [transcriptText, setTranscriptText] = useState('');

  useEffect(() => {
    // Fetch the user's progress for this course when the component mounts or the user changes.
    const fetchProgress = async () => {
      if (user) {
        try {
          // We need a new API route to proxy this request
          const res = await fetch(`/api/progress/${user.id}`);
          if (res.ok) {
            const data = await res.json();
            setCompletedLessons(new Set(data.completed_lessons));
          }
        } catch (error) {
          console.error("Failed to fetch user progress", error);
        }
      }
    };
    fetchProgress();
  }, [user]);

  const handleComplete = async (lessonId) => {
    if (!user) return;

    try {
      const res = await fetch(`/api/progress/${user.id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ lessonId }),
      });

      if (res.ok) {
        setCompletedLessons(prev => new Set(prev).add(lessonId));
      } else {
        alert('Failed to mark as complete.');
      }
    } catch (error) {
      console.error('Failed to update progress', error);
      alert('An error occurred.');
    }
  };

  // If the page is not yet generated, this will be displayed
  // initially until getStaticProps() finishes running
  if (router.isFallback) {
    return <div>Loading...</div>;
  }

  const handleGenerateQuiz = async () => {
    setIsLoadingQuiz(true);
    setShowQuizModal(true);
    try {
      // For now, we'll just use the course description to generate a quiz.
      // A real implementation would use the text content of a specific lesson.
      const res = await fetch('/api/quiz', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text_content: course.description }),
      });
      if (res.ok) {
        const data = await res.json();
        setQuizData(data);
      } else {
        setQuizData(null); // Clear old quiz data on error
      }
    } catch (error) {
      console.error('Failed to generate quiz', error);
      setQuizData(null);
    } finally {
      setIsLoadingQuiz(false);
    }
  };

  const handleShowTranscript = async (lesson) => {
    if (visibleTranscript === lesson.id) {
      setVisibleTranscript(null); // Hide if already visible
      return;
    }

    if (!lesson.transcript_url) {
      alert("Transcript not available for this lesson.");
      return;
    }

    try {
      // In a real app, the transcript_url would be a GCS URL. We can't fetch that directly
      // from the client. We'd need another API proxy route to fetch it for us.
      // For this demo, we'll just simulate the text.
      setTranscriptText("This is a placeholder for the full video transcript...");
      setVisibleTranscript(lesson.id);
    } catch (error) {
      console.error("Failed to fetch transcript", error);
      alert("Could not load transcript.");
    }
  };

  return (
    <>
      <ReportModal
        contentId={course.id}
        show={showReportModal}
        onClose={() => setShowReportModal(false)}
      />
      <Modal
        show={showQuizModal}
        onClose={() => setShowQuizModal(false)}
        title="Test Your Knowledge"
      >
        {isLoadingQuiz && <p>Generating quiz...</p>}
        {!isLoadingQuiz && quizData && <Quiz quizData={quizData} />}
        {!isLoadingQuiz && !quizData && <p>Could not load the quiz.</p>}
      </Modal>
      <div className="container">
        <Head>
          <title>{course.title}</title>
        <meta name="description" content={course.description} />
      </Head>

      <main>
        <div className={styles.titleWrapper}>
          <h1 className={styles.courseTitle}>{course.title}</h1>
          {user && <button onClick={() => setShowReportModal(true)} className={styles.reportButton}>Report</button>}
        </div>
        <p className={styles.courseDescription}>{course.description}</p>

        <div className={styles.lessonList}>
          <div className={styles.lessonHeader}>
            <h2>Lessons</h2>
            <button onClick={handleGenerateQuiz} className={styles.quizButton}>Test Your Knowledge</button>
          </div>
          <ul>
            {lessons.map((lesson) => {
              const isCompleted = completedLessons.has(lesson.id);
              return (
                <div key={lesson.id}>
                  <li className={`${styles.lessonItem} ${isCompleted ? styles.completed : ''}`}>
                    <span>{isCompleted ? '✔' : lesson.position}. {lesson.title}</span>
                    <div className={styles.lessonActions}>
                      {lesson.transcript_url && (
                        <button onClick={() => handleShowTranscript(lesson)} className={styles.transcriptButton}>
                          {visibleTranscript === lesson.id ? 'Hide' : 'Show'} Transcript
                        </button>
                      )}
                      {!isCompleted && user && (
                        <button onClick={() => handleComplete(lesson.id)} className={styles.completeButton}>
                          Mark as Complete
                        </button>
                      )}
                    </div>
                  </li>
                  {visibleTranscript === lesson.id && (
                    <div className={styles.transcriptBox}>
                      <p>{transcriptText}</p>
                    </div>
                  )}
                </div>
              );
            })}
          </ul>
        </div>

        <div className={styles.qnaSection}>
          <h2>Ask a Question</h2>
          <ChatInterface />
        </div>

        <div className={styles.reviewsSection}>
          <h2>Reviews</h2>
          <ReviewForm courseId={course.id} onReviewSubmitted={() => { /* Add refetch logic here */ }} />
          <Reviews courseId={course.id} />
        </div>

        <div className={styles.discussionSection}>
          <h2>Discussions</h2>
          <DiscussionTab courseId={course.id} />
        </div>
      </main>
    </div>
  );
}

// This function tells Next.js which dynamic paths to pre-render.
export async function getStaticPaths() {
  try {
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const res = await fetch(`${gatewayUrl}/api/content/courses`);
    const courses = await res.json();

    const paths = courses.map((course) => ({
      params: { id: course.id.toString() },
    }));

    return { paths, fallback: 'blocking' };
  } catch (error) {
    console.error('Failed to fetch paths for courses:', error);
    return { paths: [], fallback: 'blocking' };
  }
}

// This function fetches the data for a single course at build time.
export async function getStaticProps({ params }) {
  try {
    // Fetch course details and lessons from the API gateway
    const gatewayUrl = process.env.API_GATEWAY_URL || 'http://api-gateway:8080';
    const res = await fetch(`${gatewayUrl}/api/content/courses/${params.id}`);

    if (!res.ok) {
      // If the response is not ok (e.g., 404), we want to show a 404 page.
      return { notFound: true };
    }

    // Assuming the API returns an object like { course: {...}, lessons: [...] }
    const { course, lessons } = await res.json();

    return {
      props: {
        course,
        lessons,
      },
      // Re-generate the page at most once every 60 seconds
      revalidate: 60,
    };
  } catch (error) {
    console.error(`Failed to fetch data for course ${params.id}:`, error);
    // In case of an error (e.g., network issue), we can also show a 404 page.
    return { notFound: true };
  }
}
