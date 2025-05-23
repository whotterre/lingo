CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
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
    admin_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) UNIQUE NOT NULL,
    last_name VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    profile_image_url VARCHAR(255),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE languages (
    language_id SERIAL PRIMARY KEY,
    language_name VARCHAR(50) NOT NULL,
    language_code VARCHAR(5) NOT NULL, -- e.g., "yo" for Yoruba
    flag_emoji VARCHAR(10), -- 🇳🇬 for Nigerian languages
    description TEXT
);

CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    language_id INT REFERENCES languages(language_id),
    course_name VARCHAR(100) NOT NULL,
    difficulty_level VARCHAR(20) CHECK (difficulty_level IN ('Beginner', 'Intermediate', 'Advanced')),
    is_free BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lessons (
    lesson_id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(course_id),
    lesson_title VARCHAR(100) NOT NULL,
    lesson_order INT NOT NULL, -- Order in course
    xp_reward INT DEFAULT 10,
    is_unlocked BOOLEAN DEFAULT FALSE
);

CREATE TABLE exercises (
    exercise_id SERIAL PRIMARY KEY,
    lesson_id INT REFERENCES lessons(lesson_id),
    exercise_type VARCHAR(20) CHECK (exercise_type IN ('MultipleChoice', 'FillBlank', 'Speaking', 'Listening')),
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    options JSONB, -- For multiple-choice: {"options": ["Ẹ káàbọ̀", "Báwo ni?", "Ẹ ṣeun"]}
    audio_url VARCHAR(255) -- For listening exercises
);

CREATE TABLE user_progress (
    progress_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    lesson_id INT REFERENCES lessons(lesson_id),
    exercise_id INT REFERENCES exercises(exercise_id),
    is_completed BOOLEAN DEFAULT FALSE,
    score INT,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_courses (
    user_course_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id),
    course_id INT REFERENCES courses(course_id),
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completion_percentage FLOAT DEFAULT 0.0
);