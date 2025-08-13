import Head from 'next/head';
import styles from '../styles/Leaderboard.module.css';

export default function LeaderboardPage({ leaderboard, error }) {
  if (error) {
    return <div className="container"><p>{error}</p></div>;
  }

  return (
    <div className="container">
      <Head>
        <title>Leaderboard</title>
      </Head>

      <main>
        <h1 className={styles.title}>Top Learners</h1>
        <div className={styles.leaderboard}>
          <div className={styles.header}>
            <span>Rank</span>
            <span>User</span>
            <span>Score</span>
          </div>
          {leaderboard.map((entry, index) => (
            <div key={entry.user_id} className={styles.row}>
              <span className={styles.rank}>{index + 1}</span>
              <span className={styles.user}>User {entry.user_id}</span>
              <span className={styles.score}>{entry.score}</span>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}

export async function getServerSideProps() {
  try {
    // This will be a new API route we need to create
    // We use a full URL here because this runs on the server side.
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000';
    const res = await fetch(`${baseUrl}/api/leaderboard`);

    if (!res.ok) {
      throw new Error('Failed to fetch leaderboard data.');
    }

    const leaderboard = await res.json();
    return { props: { leaderboard } };
  } catch (error) {
    return { props: { error: error.message, leaderboard: [] } };
  }
}
