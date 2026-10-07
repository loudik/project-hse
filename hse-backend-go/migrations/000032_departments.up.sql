CREATE TABLE IF NOT EXISTS departments (
  id                  VARCHAR(36)  NOT NULL PRIMARY KEY,
  name                VARCHAR(255) NOT NULL UNIQUE,
  notify_on_submission BOOLEAN     NOT NULL DEFAULT FALSE, -- checked by Admin = auto-notify all members on every vessel submission
  created_at          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE users
  ADD COLUMN department_id VARCHAR(36) NULL,
  ADD CONSTRAINT fk_user_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL;