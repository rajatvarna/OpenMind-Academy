import Head from 'next/head';
import { useRouter } from 'next/router';
import styles from '../../styles/CoursePage.module.css';

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
      </main>
    </div>
  );
}

// This function tells Next.js which dynamic paths to pre-render.
export async function getStaticPaths() {
  // In a real app, you'd fetch all course IDs from the Content Service.
  const allCourseIds = [{ id: '1' }, { id: '2' }, { id: '3' }, { id: '4' }];

  const paths = allCourseIds.map((course) => ({
    params: { id: course.id.toString() },
  }));

  // { fallback: true } means other routes should be generated on-demand.
  return { paths, fallback: true };
}

// This function fetches the data for a single course at build time.
export async function getStaticProps({ params }) {
  // params contains the course `id`.
  // In a real app, you would fetch course and lesson data from your API.
  // const courseRes = await fetch(`http://content-service:3001/api/v1/courses/${params.id}`);
  // const courseData = await courseRes.json();

  const placeholderCourse = { id: params.id, title: `Course ${params.id}`, description: `This is the description for course ${params.id}.` };
  const placeholderLessons = [
    { id: 1, course_id: params.id, position: 1, title: 'Welcome to the Course' },
    { id: 2, course_id: params.id, position: 2, title: 'Core Concepts' },
    { id: 3, course_id: params.id, position: 3, title: 'Advanced Topics' },
  ];

  return {
    props: {
      course: placeholderCourse,
      lessons: placeholderLessons,
    },
    revalidate: 60,
  };
}
