CREATE TABLE user_profiles (
    user_id VARCHAR(255) PRIMARY KEY, -- Референс на users.id в другой БД
    email VARCHAR(255) UNIQUE NOT NULL, 
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    telegram VARCHAR(100),
    is_tutor BOOLEAN DEFAULT false,
    is_student BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tutor_profiles (
    user_id VARCHAR(255) PRIMARY KEY,
    specialization VARCHAR(200),
    experience_years INTEGER,
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE student_profiles (
    user_id VARCHAR(255) PRIMARY KEY,
    grade INTEGER,
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
