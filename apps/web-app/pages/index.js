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
  // In a real application, you would fetch this data from your Content Service API.
  // const res = await fetch('http://content-service:3001/api/v1/courses');
  // const courses = await res.json();

  // For now, we'll use placeholder data.
  const courses = [
    { id: 1, title: 'Introduction to Python', description: 'Learn the basics of Python programming from scratch.' },
    { id: 2, title: 'Web Development Fundamentals', description: 'Understand HTML, CSS, and JavaScript, the building blocks of the web.' },
    { id: 3, title: 'The Science of Well-Being', description: 'A course on the science behind happiness and productivity.' },
    { id: 4, title: 'Graphic Design for Beginners', description: 'Learn the core principles of graphic design.' },
  ];

  return {
    props: {
      courses,
    },
    // Optional: revalidate the data every 60 seconds
    revalidate: 60,
  };
}
