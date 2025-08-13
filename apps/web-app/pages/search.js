import Head from 'next/head';
import CourseCard from '../components/CourseCard';
import styles from '../styles/Home.module.css'; // Reuse styles from Home

export default function SearchPage({ results, query }) {
  return (
    <div className="container">
      <Head>
        <title>Search Results for "{query}"</title>
      </Head>

      <main>
        <h1>Search Results for "{query}"</h1>

        <div className={styles.grid}>
          {results && results.length > 0 ? (
            results.map((result) => (
              // Assuming search results are courses and have a compatible structure
              <CourseCard key={result.document_id} course={{ id: result.document_id, ...result.source }} />
            ))
          ) : (
            <p>No results found.</p>
          )}
        </div>
      </main>
    </div>
  );
}

export async function getServerSideProps(context) {
  const { query } = context;
  const searchQuery = query.q || '';

  if (!searchQuery) {
    return { props: { results: [], query: '' } };
  }

  try {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000';
    const res = await fetch(`${baseUrl}/api/search?q=${encodeURIComponent(searchQuery)}`);

    if (!res.ok) {
      throw new Error('Failed to fetch search results.');
    }

    const data = await res.json();
    return { props: { results: data.results, query: searchQuery } };
  } catch (error) {
    console.error(error);
    return { props: { results: [], query: searchQuery, error: error.message } };
  }
}
