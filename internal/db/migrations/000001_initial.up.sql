CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    profile_image_url VARCHAR(255),
    streak_count INT DEFAULT 0,
    xp_points INT DEFAULT 0,
    last_active_date TIMESTAMP,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE admins (
    admin_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    profile_image_url VARCHAR(255),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE languages (
    language_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    language_name VARCHAR(50) NOT NULL,
    language_code VARCHAR(5) NOT NULL,
    flag_emoji VARCHAR(10),
    description TEXT,
    CONSTRAINT unique_language_code UNIQUE (language_code)
);

CREATE TABLE courses (
    course_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    language_id UUID NOT NULL REFERENCES languages(language_id) ON DELETE CASCADE,
    description TEXT,
    course_name VARCHAR(100) NOT NULL,
    difficulty_level VARCHAR(20) CHECK (difficulty_level IN ('Beginner', 'Intermediate', 'Advanced')),
    is_free BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_course_language UNIQUE (course_name, language_id)
);

CREATE TABLE lessons (
    lesson_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    lesson_title VARCHAR(100) NOT NULL,
    lesson_order INT NOT NULL,
    xp_reward INT DEFAULT 10,
    is_unlocked BOOLEAN DEFAULT FALSE,
    CONSTRAINT unique_lesson_order_per_course UNIQUE (course_id, lesson_order)
);

CREATE TABLE exercises (
    exercise_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons(lesson_id) ON DELETE CASCADE,
    exercise_type VARCHAR(20) CHECK (exercise_type IN ('MultipleChoice', 'FillBlank', 'Speaking', 'Listening')),
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    options JSONB,
    audio_url VARCHAR(255),
    CONSTRAINT valid_multiple_choice CHECK (
        exercise_type != 'MultipleChoice' OR 
        (options IS NOT NULL AND jsonb_typeof(options->'options') = 'array')
    )
);

CREATE TABLE user_progress (
    progress_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL REFERENCES lessons(lesson_id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(exercise_id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE,
    score INT,
    completed_at TIMESTAMP,
    CONSTRAINT unique_user_exercise UNIQUE (user_id, exercise_id),
    CONSTRAINT valid_score CHECK (score IS NULL OR (score >= 0 AND score <= 100))
);

CREATE TABLE user_courses (
    user_course_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completion_percentage FLOAT DEFAULT 0.0,
    CONSTRAINT unique_user_course UNIQUE (user_id, course_id),
    CONSTRAINT valid_percentage CHECK (completion_percentage >= 0.0 AND completion_percentage <= 100.0)
);