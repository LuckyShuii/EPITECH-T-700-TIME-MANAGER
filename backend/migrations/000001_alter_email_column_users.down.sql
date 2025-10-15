-- TEST migration to alter the email column in users table to VARCHAR(321)
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(320);