ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

CREATE TABLE IF NOT EXISTS email_verifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_email_verifications_token_hash ON email_verifications(token_hash);
