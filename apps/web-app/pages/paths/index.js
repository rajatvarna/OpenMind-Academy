import Head from 'next/head';
import Link from 'next/link';
import styles from '../../styles/Home.module.css'; // Reuse styles

export default function PathsPage({ paths }) {
  return (
    <div className="container">
      <Head>
        <title>Learning Paths</title>
      </Head>
      <main>
        <h1 className={styles.title}>Learning Paths</h1>
        <p className={styles.description}>
          Follow our curated paths to master a new skill from start to finish.
        </p>
        <div className={styles.grid}>
          {paths.map((path) => (
            <Link key={path.id} href={`/paths/${path.id}`} legacyBehavior>
              <a className={styles.card}>
                <h3>{path.title} &rarr;</h3>
                <p>{path.description}</p>
              </a>
            </Link>
          ))}
        </div>
      </main>
    </div>
  );
}

export async function getStaticProps() {
  // In a real app, we would fetch this from /api/paths
  const placeholderPaths = [
    { id: 1, title: 'Web Development Career Path', description: 'Everything you need to become a web developer.' },
    { id: 2, title: 'Data Science A-Z', description: 'Learn Python, Pandas, and Scikit-learn.' },
  ];

  return {
    props: {
      paths: placeholderPaths,
    },
    revalidate: 60,
  };
}
