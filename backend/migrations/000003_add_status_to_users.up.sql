-- EITHER active, disabled or pending (if user has not created his password yet)
ALTER TABLE users
ADD COLUMN status VARCHAR(15) NOT NULL DEFAULT 'pending';