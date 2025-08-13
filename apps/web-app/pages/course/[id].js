import Head from 'next/head';
import { useRouter } from 'next/router';
import styles from '../../styles/CoursePage.module.css';
import ChatInterface from '../../components/ChatInterface';

export default function CoursePage({ course, lessons }) {
  const router = useRouter();

  // If the page is not yet generated, this will be displayed
  // initially until getStaticProps() finishes running
  if (router.isFallback) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container">
      <Head>
        <title>{course.title}</title>
        <meta name="description" content={course.description} />
      </Head>

      <main>
        <h1 className={styles.courseTitle}>{course.title}</h1>
        <p className={styles.courseDescription}>{course.description}</p>

        <div className={styles.lessonList}>
          <h2>Lessons</h2>
          <ul>
            {lessons.map((lesson) => (
              <li key={lesson.id} className={styles.lessonItem}>
                <span>{lesson.position}. {lesson.title}</span>
              </li>
            ))}
          </ul>
        </div>

        <div className={styles.qnaSection}>
          <h2>Ask a Question</h2>
          <ChatInterface />
        </div>
      </main>
    </div>
  );
}

// This function tells Next.js which dynamic paths to pre-render.
export async function getStaticPaths() {
  try {
    const res = await fetch('http://api-gateway:8080/api/content/courses');
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
    const res = await fetch(`http://api-gateway:8080/api/content/courses/${params.id}`);

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
