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

-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE
    admin_id = $1 RETURNING admin_id,
    first_name,
    last_name,
    email,
    profile_image_url,
    joined_at;

-- name: GetAllAdmins :many
SELECT
    admin_id,
    first_name,
    last_name,
    email,
    profile_image_url,
    joined_at
FROM admins
ORDER BY joined_at DESC
LIMIT $1
OFFSET
    $2;

-- Maintenance Functionality --
-- Creational Queries

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
        language_id,
        difficulty_level,
        is_free
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateLesson :one
INSERT INTO
    lessons (
        lesson_title,
        lesson_order,
        xp_reward,
        is_unlocked
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateExercise :one
INSERT INTO
    exercises (
        exercise_type,
        question_text,
        correct_answer,
        options,
        audio_url
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

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