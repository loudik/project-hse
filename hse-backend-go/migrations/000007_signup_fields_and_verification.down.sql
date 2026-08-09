DROP TABLE IF EXISTS email_verifications;

ALTER TABLE users
  DROP COLUMN status,
  DROP COLUMN accepted_terms_at,
  DROP COLUMN website,
  DROP COLUMN location,
  DROP COLUMN organization_name,
  DROP COLUMN position,
  DROP COLUMN phone_number;
