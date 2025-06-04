-- name: GetAdminForLogin :one
SELECT admin_id, email, password
FROM admins
WHERE
    email = $1
LIMIT 1;

-- name: CreateAdmin :one
INSERT INTO
    admins (
        first_name,
        last_name,
        email,
        password,
        profile_image_url
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAdminById :one
SELECT
    admin_id,
    first_name,
    last_name,
    email,
    profile_image_url,
    joined_at
FROM admins
WHERE
    admin_id = $1
LIMIT 1;

-- name: GetAdminByEmail :one
SELECT
    admin_id,
    first_name,
    last_name,
    email,
    profile_image_url,
    joined_at
FROM admins
WHERE
    email = $1
LIMIT 1;

-- name: UpdateAdmin :exec
UPDATE admins
SET
    first_name = $1,
    last_name = $2,
    email = $3,
    password = $4,
    profile_image_url = $5
WHERE
    admin_id = $6 RETURNING admin_id,
    first_name,
    last_name,
    email,
    profile_image_url,
    joined_at;

-- Delete admin by ID
-- name: DeleteAdmin :exec
DELETE FROM admins WHERE admin_id = $1;

-- name: CreateLanguage :one
INSERT INTO
    languages (
        language_name,
        language_code,
        flag_emoji,
        description
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateCourse :one
INSERT INTO
    courses (
        course_name,
        description,
        language_id,
        difficulty_level,
        is_free
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateLesson :one
INSERT INTO
    lessons (
        lesson_title,
        course_id,
        lesson_order,
        xp_reward,
        is_unlocked
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateExercise :one
INSERT INTO
    exercises (
        lesson_id,
        exercise_type,
        question_text,
        correct_answer,
        options,
        audio_url
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: CreateUserProgress :one
INSERT INTO
    user_progress (
        user_id,
        lesson_id,
        score,
        completed_at
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateUserCourse :one
INSERT INTO
    user_courses (
        user_id,
        course_id,
        enrollment_date,
        completion_percentage
    )
VALUES (
        $1,
        $2,
        CURRENT_TIMESTAMP,
        0.0
    ) RETURNING *;

-- Retrieval Queries

-- name: GetAllLanguages :many
SELECT
    language_id,
    language_name,
    language_code,
    flag_emoji,
    description
FROM languages
ORDER BY language_name ASC
LIMIT $1
OFFSET
    $2;

-- name: GetLanguageById :one
SELECT
    language_id,
    language_name,
    language_code,
    flag_emoji,
    description
FROM languages
WHERE
    language_id = $1
LIMIT 1;

-- name: GetLanguageByName :one
SELECT
    language_id,
    language_name,
    language_code,
    flag_emoji,
    description
FROM languages
WHERE
    language_name = $1
LIMIT 1;

-- name: GetAllCourses :many
SELECT * FROM courses;

-- name: GetAllCoursesByLanguage :many
SELECT * FROM courses WHERE language_id = $1;
-- name: GetAllExercises :many
SELECT
    exercise_id,
    lesson_id,
    exercise_type,
    question_text,
    correct_answer,
    options,
    audio_url
FROM exercises
LIMIT $1
OFFSET
    $2;

-- name: GetExerciseById :one
SELECT
    exercise_id,
    lesson_id,
    exercise_type,
    question_text,
    correct_answer,
    options,
    audio_url
FROM exercises
WHERE
    exercise_id = $1
LIMIT 1;

-- name: GetAllLessons :many
SELECT * FROM lessons LIMIT $1 OFFSET $2;

-- name: GetLessonById :one
SELECT
    lesson_id,
    lesson_title,
    course_id,
    lesson_order,
    xp_reward,
    is_unlocked
FROM lessons
WHERE
    lesson_id = $1
LIMIT 1;

-- name: GetLessonsByCourseId :many
SELECT
    lesson_id,
    lesson_title,
    course_id,
    lesson_order,
    xp_reward,
    is_unlocked
FROM lessons
WHERE
    course_id = $1
ORDER BY lesson_order ASC
LIMIT $2
OFFSET
    $3;

-- name: GetUserProgressByUserId :many
SELECT up.progress_id, up.user_id, up.lesson_id, up.score, up.completed_at, l.lesson_title
FROM user_progress up
    JOIN lessons l ON up.lesson_id = l.lesson_id
WHERE
    up.user_id = $1
ORDER BY up.completed_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetExercisesByLessonId :many
SELECT
    exercise_id,
    lesson_id,
    exercise_type,
    question_text,
    correct_answer,
    options,
    audio_url
FROM exercises
WHERE
    lesson_id = $1
ORDER BY exercise_id ASC
LIMIT $2
OFFSET
    $3;


-- name: UpdateAdminDetails :exec
UPDATE admins
SET
    first_name = $1,
    last_name = $2,
    email = $3,
    profile_image_url = $4
WHERE
    admin_id = $5;

-- Update admin password
-- name: UpdateAdminPassword :exec
UPDATE admins SET password = $1 WHERE admin_id = $2;

-- Update language details
-- name: UpdateLanguageDetails :exec
UPDATE languages
SET
    language_name = $1,
    language_code = $2,
    flag_emoji = $3,
    description = $4
WHERE
    language_id = $5;

-- Update course details
-- name: UpdateCourseDetails :exec
UPDATE courses
SET
    course_name = $1,
    description = $2,
    difficulty_level = $3,
    is_free = $4
WHERE
    course_id = $5;

-- Update lesson details
-- name: UpdateLessonDetails :exec
UPDATE lessons
SET
    lesson_title = $1,
    lesson_order = $2,
    xp_reward = $3,
    is_unlocked = $4
WHERE
    lesson_id = $5;

-- Update exercise details
-- name: UpdateExerciseDetails :exec
UPDATE exercises
SET
    exercise_type = $1,
    question_text = $2,
    correct_answer = $3,
    options = $4,
    audio_url = $5
WHERE
    exercise_id = $6;

-- Update user progress
-- name: UpdateUserProgress :exec
UPDATE user_progress
SET
    is_completed = $1,
    score = $2,
    completed_at = $3
WHERE
    progress_id = $4;

-- Update user course progress
-- name: UpdateUserCourseProgress :exec
UPDATE user_courses
SET
    completion_percentage = $1
WHERE
    user_course_id = $2;

-- Delete user by ID
-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- Delete user progress by user ID
-- name: DeleteUserProgressByUserId :exec
DELETE FROM user_progress WHERE user_id = $1;

-- Delete user progress by lesson ID
-- name: DeleteUserProgressByLessonId :exec
DELETE FROM user_progress WHERE lesson_id = $1;

-- Delete user progress by exercise ID
-- name: DeleteUserProgressByExerciseId :exec
DELETE FROM user_progress WHERE exercise_id = $1;

-- Delete user courses by user ID
-- name: DeleteUserCoursesByUserId :exec
DELETE FROM user_courses WHERE user_id = $1;

-- Delete user courses by course ID
-- name: DeleteUserCoursesByCourseId :exec
DELETE FROM user_courses WHERE course_id = $1;

-- Delete course by ID
-- name: DeleteCourse :exec
DELETE FROM courses WHERE course_id = $1;

-- Delete lesson by ID
-- name: DeleteLesson :exec
DELETE FROM lessons WHERE lesson_id = $1;

-- Delete lessons by course ID
-- name: DeleteLessonsByCourseId :exec
DELETE FROM lessons WHERE course_id = $1;

-- Delete exercise by ID
-- name: DeleteExercise :exec
DELETE FROM exercises WHERE exercise_id = $1;

-- Delete exercises by lesson ID
-- name: DeleteExercisesByLessonId :exec
DELETE FROM exercises WHERE lesson_id = $1;

-- Delete language by ID
-- name: DeleteLanguage :exec
DELETE FROM languages WHERE language_id = $1;

-- Delete courses by language ID
-- name: DeleteCoursesByLanguageId :exec
DELETE FROM courses WHERE language_id = $1;