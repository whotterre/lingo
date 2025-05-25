# Lingo

Lingo is a Nigerian Language Learning App designed to help users learn various Nigerian languages through structured courses, lessons, and exercises. The app provides an admin interface for managing languages, courses, lessons, and exercises, as well as user progress tracking.

---

## Features

### Admin Features
- **Authentication**: Admins can sign up, log in, and manage their accounts.
- **Language Management**: Create, update, delete, and retrieve languages.
- **Course Management**: Create, update, delete, and retrieve courses associated with specific languages.
- **Lesson Management**: Create, update, delete, and retrieve lessons associated with specific courses.
- **Exercise Management**: Create, update, delete, and retrieve exercises associated with specific lessons.

### User Features
- **Authentication**: Users can sign up, log in, and manage their accounts.
- **Progress Tracking**: Track user progress through lessons and exercises.
- **Course Enrollment**: Users can enroll in courses and track their completion percentage.

---

## Project Structure

### Handlers
The `handlers` package contains the logic for handling HTTP requests. Key handlers include:
- `AdminHandler`: Handles admin-related operations such as managing languages, courses, lessons, and exercises.

### Database
The `db` package contains SQL queries and database interaction logic. It uses the `sqlc` tool to generate type-safe database access code.

### Routes
The `server.go` file defines the API routes for the application. Routes are grouped into:
- **Authentication Routes**: For admin and user authentication.
- **Admin Routes**: For managing admin details.
- **Language Routes**: For managing languages.
- **Course Routes**: For managing courses.
- **Lesson Routes**: For managing lessons.
- **Exercise Routes**: For managing exercises.
- **User Routes**: For user-specific operations.

---

## API Endpoints

### Authentication Routes
- `POST /auth/learner/signup`: User signup.
- `POST /auth/learner/login`: User login.
- `POST /auth/learner/refresh`: Refresh user token.
- `POST /auth/admin/signup`: Admin signup.
- `POST /auth/admin/login`: Admin login.
- `POST /auth/admin/refresh`: Refresh admin token.

### Admin Routes
- `PUT /admin/details/:adminId`: Update admin details.
- `PUT /admin/password/:adminId`: Update admin password.

### Language Routes
- `POST /admin/language/create`: Create a new language.
- `PUT /admin/language/:languageId`: Update a language.
- `DELETE /admin/language/:languageId`: Delete a language.
- `GET /admin/lesson/languages/all`: Retrieve all languages.

### Course Routes
- `POST /admin/course/create/:langId`: Create a new course.
- `PUT /admin/course/:courseId`: Update a course.
- `DELETE /admin/course/:courseId`: Delete a course.
- `GET /admin/lesson/courses/all`: Retrieve all courses.

### Lesson Routes
- `POST /admin/lesson/create/:courseId`: Create a new lesson.
- `PUT /admin/lesson/:lessonId`: Update a lesson.
- `DELETE /admin/lesson/:lessonId`: Delete a lesson.
- `GET /admin/lesson/lessons/all`: Retrieve all lessons.
- `GET /admin/lesson/lessons/by-course/:courseId`: Retrieve lessons by course.

### Exercise Routes
- `POST /admin/exercise/create`: Create a new exercise.
- `PUT /admin/exercise/:exerciseId`: Update an exercise.
- `DELETE /admin/exercise/:exerciseId`: Delete an exercise.
- `GET /admin/exercise/:exerciseId`: Retrieve an exercise by ID.
- `GET /admin/exercise/exercises/all`: Retrieve all exercises.
- `GET /admin/exercise/exercises/by-lesson/:lessonId`: Retrieve exercises by lesson.

### User Routes
- `GET /users/me`: Retrieve current user details.
- `PUT /users/me`: Update current user details.
- `GET /users/id`: Retrieve user by ID.

---

## Database Schema

### Tables
- **Admins**: Stores admin details.
- **Users**: Stores user details.
- **Languages**: Stores language details.
- **Courses**: Stores course details associated with languages.
- **Lessons**: Stores lesson details associated with courses.
- **Exercises**: Stores exercise details associated with lessons.
- **User Progress**: Tracks user progress in lessons and exercises.
- **User Courses**: Tracks user enrollment and completion percentage in courses.

---

## Setup Instructions

### Prerequisites
- Go 1.18+
- PostgreSQL
- `sqlc` for database query generation

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/whotterre/lingo.git
   cd lingo
2. Install necessary dependencies
   ```bash
   go mod tidy
3. Run the server code 
   If you have GNU make installed
   ```bash
   make run
   ```
   Otherwise, 
   ```bash
   cd cmd/api && go run .
