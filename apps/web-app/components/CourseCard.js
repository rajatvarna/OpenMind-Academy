import React from 'react';
import Link from 'next/link';

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
      <a className="m-4 p-6 text-left text-inherit no-underline border border-gray-200 rounded-lg transition-colors duration-150 ease-in-out w-5/12 hover:text-blue-600 hover:border-blue-600 focus:text-blue-600 focus:border-blue-600 active:text-blue-600 active:border-blue-600">
        <h3 className="mb-4 text-2xl">{course.title} &rarr;</h3>
        <p className="m-0 text-xl leading-normal">{course.description}</p>
      </a>
    </Link>
  );
};

export default CourseCard;
