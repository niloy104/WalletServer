-- Delete a specific user (id = 5)
DELETE FROM users WHERE id = 5;

-- Delete all users (keep table structure)
DELETE FROM users;

-- Delete all users and reset auto-increment (PostgreSQL)
TRUNCATE TABLE users RESTART IDENTITY;