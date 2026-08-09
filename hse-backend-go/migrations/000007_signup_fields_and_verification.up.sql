-- ============ Extra profile fields on users (US 1.1) ============
ALTER TABLE users
  ADD COLUMN phone_number      VARCHAR(50)  NULL AFTER email,
  ADD COLUMN position          VARCHAR(150) NULL AFTER phone_number,
  ADD COLUMN organization_name VARCHAR(255) NULL AFTER position,
  ADD COLUMN location          VARCHAR(255) NULL AFTER organization_name,
  ADD COLUMN website           VARCHAR(255) NULL AFTER location,
  ADD COLUMN accepted_terms_at DATETIME     NULL AFTER website,
  ADD COLUMN status             VARCHAR(20) NOT NULL DEFAULT 'Pending' AFTER organization_id; -- Pending | Active | Suspended

-- Existing users (seeded admin, anyone created before this migration)
-- should stay usable - mark them Active.
UPDATE users SET status = 'Active';

-- ============ Email verification tokens ============
CREATE TABLE IF NOT EXISTS email_verifications (
  id         VARCHAR(36)  NOT NULL PRIMARY KEY,
  user_id    VARCHAR(36)  NOT NULL,
  token      VARCHAR(100) NOT NULL UNIQUE,
  expires_at DATETIME     NOT NULL,
  used_at    DATETIME     NULL,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
