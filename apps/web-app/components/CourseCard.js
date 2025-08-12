import React from 'react';
import Link from 'next/link';
import styles from '../styles/CourseCard.module.css';

const CourseCard = ({ course }) => {
  return (
    <Link href={`/course/${course.id}`} legacyBehavior>
      <a className={styles.card}>
        <h3>{course.title} &rarr;</h3>
        <p>{course.description}</p>
      </a>
    </Link>
  );
};

export default CourseCard;
