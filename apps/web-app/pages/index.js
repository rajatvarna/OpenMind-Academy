import Head from 'next/head';
import CourseCard from '../components/CourseCard';
import styles from '../styles/Home.module.css';

export default function Home({ courses }) {
  return (
    <div className="container">
      <Head>
        <title>Free Education Platform - Courses</title>
        <meta name="description" content="Browse our list of free courses." />
      </Head>

      <main>
        <h1 className={styles.title}>
          Welcome to the Future of Learning
        </h1>

        <p className={styles.description}>
          Explore our community-generated courses on any topic imaginable.
        </p>

        <div className={styles.grid}>
          {courses.map((course) => (
            <CourseCard key={course.id} course={course} />
          ))}
        </div>
      </main>
    </div>
  );
}

// This function runs at build time on the server.
export async function getStaticProps() {
  let courses = [];
  try {
    // Fetch data from the API gateway, which routes to the Content Service.
    // In a real K8s setup, this would be the internal service name.
    // For local dev, it might be http://localhost:8080/api/content/courses
    const res = await fetch('http://api-gateway:8080/api/content/courses');

    if (res.ok) {
      courses = await res.json();
    } else {
      // Log an error to the server-side console
      console.error('Failed to fetch courses:', res.status, res.statusText);
    }
  } catch (error) {
    console.error('An error occurred while fetching courses:', error);
  }

  // The page will be rendered with the fetched courses, or an empty array if the fetch failed.
  return {
    props: {
      courses,
    },
    // Re-generate the page at most once every 60 seconds
    revalidate: 60,
  };
}
