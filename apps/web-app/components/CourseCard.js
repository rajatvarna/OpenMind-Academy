import React from 'react';
import Link from 'next/link';
import styles from '../styles/CourseCard.module.css';

/**
 * A card component that displays a summary of a course.
 * It links to the full course page.
 *
 * @param {object} course - The course object to display.
 * @param {string} course.id - The unique identifier for the course.
 * @param {string} course.title - The title of the course.
 * @param {string} course.description - A short description of the course.
 */
const CourseCard = ({ course }) => {
  return (
    // The entire card is a link to the course's detail page.
    <Link href={`/course/${course.id}`} legacyBehavior>
      <a className={styles.card}>
        <h3>{course.title} &rarr;</h3>
        <p>{course.description}</p>
      </a>
    </Link>
  );
};

export default CourseCard;
