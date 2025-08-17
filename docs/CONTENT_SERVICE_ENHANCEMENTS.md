# Proposed Feature Enhancements for Content Service

Based on a review of the existing `content-service`, here are five proposed feature enhancements that would significantly improve its functionality and provide more value to course creators and students.

---

### 1. Course Drafts and Publishing Workflow

**Description:**
Currently, courses are implicitly "live" as soon as they are created. A more robust system would allow authors to save courses as drafts. This would involve adding a `status` field to the `courses` table (e.g., with values like `draft`, `in_review`, `published`, `archived`). This allows authors to build their courses over time and only make them public when they are ready.

**Components Affected:**
- `content-service`:
    - `courses` model and table would need a `status` column.
    - API endpoints would need to be updated to filter based on status (e.g., `GetAllCoursesHandler` should only return `published` courses).
    - New endpoints would be needed to manage the status transitions (e.g., `POST /courses/:courseId/publish`).
- `web-app`: UI changes to show the status of a course and provide buttons for authors to publish or archive their courses.

---

### 2. Versioning for Lessons

**Description:**
To prevent accidental data loss and allow for better content management, a versioning system for lessons could be implemented. Whenever a lesson's content is updated, a new version would be created instead of overwriting the existing one. This would allow authors to view a history of changes and revert to a previous version if needed.

**Components Affected:**
- `content-service`:
    - A new `lesson_versions` table would be needed to store the history of lesson content.
    - The `UpdateLesson` logic would need to be changed to create a new version instead of updating in-place.
    - New endpoints to view version history and revert to a specific version.

---

### 3. Course Categories and Tags

**Description:**
To improve the discoverability of courses, a system for categorizing and tagging them would be beneficial. Authors could assign one or more categories (e.g., "Web Development", "Data Science") and free-form tags (e.g., "React", "Python", "Beginner") to their courses. Users could then filter and search for courses based on these taxonomies.

**Components Affected:**
- `content-service`:
    - New tables would be required: `categories`, `tags`, `course_categories`, and `course_tags`.
    - The `CreateCourse` and `UpdateCourse` handlers would need to be updated to handle category and tag associations.
    - The `GetAllCoursesHandler` would need to support filtering by category or tag.
- `search-service`: The search index would need to be updated to include categories and tags.
- `web-app`: UI for authors to add tags/categories and for users to filter by them.

---

### 4. Enhanced Video Content Support

**Description:**
The current `Lesson` model has a simple `VideoURL` field. This could be expanded to provide a richer video experience. This would involve adding fields for video duration, subtitles/captions (in multiple languages), and potentially different video quality levels.

**Components Affected:**
- `content-service`: The `lessons` model would be updated with new fields (e.g., `duration_seconds`, `subtitles` as a JSONB field).
- `video-processing-service`: This service could be enhanced to automatically generate transcripts and different video resolutions upon upload.
- `web-app` / `mobile-app`: The video player component would be updated to use the new fields to display duration and provide options for subtitles.

---

### 5. Analytics for Course Authors

**Description:**
Empower course creators by providing them with analytics for their courses. A new endpoint could provide key metrics such as the total number of student enrollments, the course completion rate, average rating over time, and engagement metrics for each lesson.

**Components Affected:**
- `content-service`: A new endpoint (e.g., `GET /courses/:courseId/analytics`) would be needed. This handler would likely need to query data from multiple services.
- `user-service`: Would need to provide data on user progress for a given course.
- `gamification-service`: Could provide data on quiz attempts and scores.
- `web-app`: A new "Dashboard" or "Analytics" page for course authors to view these metrics.
