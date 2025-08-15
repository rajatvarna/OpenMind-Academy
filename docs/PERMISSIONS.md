# Role-Based Access Control (RBAC) Permissions

This document outlines the permissions for each user role in the system. The API gateway is responsible for enforcing these rules.

## Roles

- **`user`**: A regular authenticated user.
- **`moderator`**: A user with additional privileges to manage content and users.
- **`admin`**: A user with full access to the system.

## Permission Matrix

| Route                                  | Method | `user`                                 | `moderator` | `admin` | Notes                               |
| -------------------------------------- | ------ | -------------------------------------- | ----------- | ------- | ----------------------------------- |
| `/api/users/register`                  | `POST` | Public                                 | Public      | Public  |                                     |
| `/api/users/login`                     | `POST` | Public                                 | Public      | Public  |                                     |
| `/api/users/profile`                   | `GET`  | Own                                    | Own         | Own     | A user can only get their own profile. |
| `/api/users/:userId/progress`          | `GET`  | Own                                    | All         | All     | Moderators/admins can view any user's progress. |
| `/api/users/:userId/progress`          | `POST` | Own                                    | No          | No      | Only users can update their own progress. |
| `/api/users/:userId/full-profile`      | `GET`  | Own                                    | All         | All     | Moderators/admins can view any user's profile. |
| `/api/content/courses`                 | `GET`  | Public                                 | Public      | Public  |                                     |
| `/api/content/courses/featured`        | `GET`  | Public                                 | Public      | Public  |                                     |
| `/api/content/courses/:courseId`       | `GET`  | Public                                 | Public      | Public  |                                     |
| `/api/content/courses/:courseId/reviews`| `GET`  | Public                                 | Public      | Public  |                                     |
| `/api/content/reviews`                 | `POST` | Yes                                    | Yes         | Yes     | Any authenticated user can post a review. |
| `/api/ugc/submit`                      | `POST` | Yes                                    | Yes         | Yes     | Any authenticated user can submit content. |
| `/api/ugc/report`                      | `POST` | Yes                                    | Yes         | Yes     | Any authenticated user can report content. |
| `/api/qna/query`                       | `POST` | Yes                                    | Yes         | Yes     |                                     |
| `/api/qna/generate-quiz`               | `POST` | Yes                                    | Yes         | Yes     |                                     |
| `/api/gamification/leaderboard`        | `GET`  | Public                                 | Public      | Public  |                                     |
| `/api/gamification/users/:userId/stats`| `GET`  | Own                                    | All         | All     | Moderators/admins can view any user's stats. |

### Key:
- **Public**: No authentication required.
- **Yes**: Authenticated users with this role have access.
- **Own**: The user can only access the resource if the `:userId` in the path matches their own user ID.
- **All**: The user can access the resource for any `:userId`.
- **No**: The user does not have access.
