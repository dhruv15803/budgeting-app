DROP INDEX IF EXISTS idx_email_verifications_token_hash;
DROP TABLE IF EXISTS email_verifications;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_unique;
