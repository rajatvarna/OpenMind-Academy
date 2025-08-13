import Head from 'next/head';
import Link from 'next/link';
import styles from '../../styles/PathPage.module.css';

export default function PathDetailPage({ path }) {
  if (!path) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container">
      <Head>
        <title>{path.title}</title>
      </Head>
      <main>
        <h1 className={styles.pathTitle}>{path.title}</h1>
        <p className={styles.pathDescription}>{path.description}</p>
        <div className={styles.courseList}>
          {path.courses.map((course, index) => (
            <div key={course.id} className={styles.courseItem}>
              <span className={styles.stepNumber}>{index + 1}</span>
              <div className={styles.courseInfo}>
                <Link href={`/course/${course.id}`}>
                  <a className={styles.courseTitle}>{course.title}</a>
                </Link>
                <p className={styles.courseDescription}>{course.description}</p>
              </div>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}

export async function getStaticPaths() {
  // Fetch all path IDs
  return {
    paths: [{ params: { id: '1' } }, { params: { id: '2' } }],
    fallback: 'blocking',
  };
}

export async function getStaticProps({ params }) {
  // Fetch data for a single path
  const placeholderPath = {
    id: params.id,
    title: `Learning Path ${params.id}`,
    description: 'A detailed description of this learning path.',
    courses: [
      { id: 1, title: 'Introduction to Python', description: 'Learn the basics.' },
      { id: 2, title: 'Web Development Fundamentals', description: 'Learn HTML, CSS, JS.' },
    ],
  };

  return {
    props: {
      path: placeholderPath,
    },
    revalidate: 60,
  };
}
