CREATE TABLE IF NOT EXISTS notifications (
  id          VARCHAR(36)   NOT NULL PRIMARY KEY,
  user_id     VARCHAR(36)   NOT NULL,

  title       VARCHAR(255)  NOT NULL,
  message     VARCHAR(1000) NULL,
  link        VARCHAR(500)  NULL,   -- frontend route to open, e.g. /vessel/review/{id}

  is_read     BOOLEAN       NOT NULL DEFAULT FALSE,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_notifications_user (user_id, is_read),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;