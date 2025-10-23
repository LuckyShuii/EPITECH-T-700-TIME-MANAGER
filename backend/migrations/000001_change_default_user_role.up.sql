-- TEST migration to alter the roles column in users table to change default role: instead of 'user' to 'employee'
ALTER TABLE users ALTER COLUMN roles SET DEFAULT '{employee}';