-- Best effort: fails if any row already has password_hash = NULL
-- (e.g. accounts created via Microsoft OAuth). Clean those up manually
-- before rolling back this far if needed.
ALTER TABLE users MODIFY password_hash VARCHAR(255) NOT NULL;
