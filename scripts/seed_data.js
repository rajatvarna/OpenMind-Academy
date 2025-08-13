// This script is for seeding the database with initial course and lesson data.
// It simulates the actions a user or admin would take to create content via the API.

// In a real project, you might use a library like 'node-fetch' if running this
// outside a browser environment. For this example, we'll just define the logic.

const API_GATEWAY_URL = 'http://localhost:8080'; // Assuming local dev setup

const COURSES_TO_SEED = [
  {
    title: 'Introduction to Programming with Python',
    description: 'Learn the fundamentals of programming using the Python language. No prior experience required.',
    authorId: 1, // Assuming an admin or system user with ID 1
    lessons: [
      { title: 'Getting Started', textContent: 'This lesson covers setting up your Python environment.' },
      { title: 'Variables and Data Types', textContent: 'Learn about the basic data types in Python.' },
      { title: 'Your First Program', textContent: 'Write and run your first "Hello, World!" program.' },
    ],
  },
  {
    title: 'Fundamentals of Graphic Design',
    description: 'Explore the core principles of graphic design, including color theory, typography, and layout.',
    authorId: 1,
    lessons: [
      { title: 'The Elements of Design', textContent: 'Learn about line, shape, form, and texture.' },
      { title: 'Understanding Color', textContent: 'An introduction to the color wheel and color harmony.' },
    ],
  },
];

async function seedData() {
  console.log('Starting data seeding...');

  for (const courseData of COURSES_TO_SEED) {
    try {
      // 1. Create the course
      console.log(`Creating course: "${courseData.title}"`);
      const courseRes = await fetch(`${API_GATEWAY_URL}/api/content/courses`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: courseData.title,
          description: courseData.description,
          author_id: courseData.authorId,
        }),
      });
      const newCourse = await courseRes.json();
      console.log(`  -> Course created with ID: ${newCourse.id}`);

      // 2. Create lessons for the course
      for (const lessonData of courseData.lessons) {
        console.log(`  Creating lesson: "${lessonData.title}"`);
        const lessonRes = await fetch(`${API_GATEWAY_URL}/api/content/lessons`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            title: lessonData.title,
            text_content: 'Initial text content placeholder.', // The real text is in the UGC submission
            course_id: newCourse.id,
            position: courseData.lessons.indexOf(lessonData) + 1,
          }),
        });
        const newLesson = await lessonRes.json();
        console.log(`    -> Lesson created with ID: ${newLesson.id}`);

        // 3. Submit the lesson text for video generation
        console.log(`    Submitting UGC for lesson ID: ${newLesson.id}`);
        await fetch(`${API_GATEWAY_URL}/api/ugc/submit`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                lessonId: newLesson.id,
                textContent: lessonData.textContent,
            }),
        });
        console.log(`      -> UGC submitted successfully.`);
      }
    } catch (error) {
      console.error(`Failed to seed course "${courseData.title}":`, error);
    }
  }

  console.log('Data seeding complete.');
}

// To run this script, you would execute `node scripts/seed_data.js` in your terminal.
// We are just defining it here as we can't run it in this environment.
// seedData();
