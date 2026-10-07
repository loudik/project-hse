CREATE TABLE IF NOT EXISTS vessel_review_requests (
  id                 VARCHAR(36) NOT NULL PRIMARY KEY,
  application_id     VARCHAR(36) NOT NULL,
  requested_staff_id VARCHAR(36) NOT NULL,
  requested_by       VARCHAR(36) NOT NULL,
  created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (application_id) REFERENCES vessel_applications(id) ON DELETE CASCADE,
  FOREIGN KEY (requested_staff_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (requested_by) REFERENCES users(id) ON DELETE RESTRICT,
  UNIQUE KEY unique_review_request (application_id, requested_staff_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Once a Staff member is granted access to review an application (via the
-- row above), they can view its full detail and leave comments - but ONLY
-- for that specific application, never a general browse list.
CREATE TABLE IF NOT EXISTS vessel_review_comments (
  id             VARCHAR(36) NOT NULL PRIMARY KEY,
  application_id VARCHAR(36) NOT NULL,
  author_id      VARCHAR(36) NOT NULL,
  comment        TEXT NOT NULL,
  created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (application_id) REFERENCES vessel_applications(id) ON DELETE CASCADE,
  FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE RESTRICT,
  INDEX idx_comments_application (application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;