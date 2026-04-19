

CREATE TYPE user_role AS ENUM ('user','admin');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    username TEXT UNIQUE,
    password TEXT NOT NULL,
    image_url TEXT,
    role user_role NOT NULL DEFAULT 'user',
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);