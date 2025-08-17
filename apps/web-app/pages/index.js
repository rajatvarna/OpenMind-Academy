import Head from 'next/head';
import CourseCard from '../components/CourseCard';

export default function Home({ courses, featuredCourses }) {
  return (
    <div className="container mx-auto px-4">
      <Head>
        <title>Free Education Platform - Courses</title>
        <meta name="description" content="Browse our list of free courses." />
      </Head>

      <main className="py-10">
        <h1 className="text-5xl font-bold text-center">
          Welcome to the Future of Learning
        </h1>

        <p className="text-xl text-center mt-4 mb-12">
          Explore our community-generated courses on any topic imaginable.
        </p>

        {featuredCourses && featuredCourses.length > 0 && (
          <div className="mb-16">
            <h2 className="text-3xl font-bold text-center mb-8">Featured Courses</h2>
            <div className="flex flex-wrap justify-center">
              {featuredCourses.map((course) => (
                <CourseCard key={course.id} course={course} />
              ))}
            </div>
          </div>
        )}

        <h2 className="text-3xl font-bold text-center mb-8">All Courses</h2>
        <div className="flex flex-wrap justify-center">
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
  // Use a full URL for server-side fetching
  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000';

  // Fetch all courses and featured courses in parallel
  const [coursesRes, featuredRes] = await Promise.all([
    fetch(`${baseUrl}/api/courses`),
    fetch(`${baseUrl}/api/courses/featured`)
  ]);

  const courses = coursesRes.ok ? await coursesRes.json() : [];
  const featuredCourses = featuredRes.ok ? await featuredRes.json() : [];

  return {
    props: {
      courses,
      featuredCourses,
    },
    // Re-generate the page at most once every 60 seconds
    revalidate: 60,
  };
}
