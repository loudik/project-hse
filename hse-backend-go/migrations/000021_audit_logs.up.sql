CREATE TABLE IF NOT EXISTS audit_logs (
  id            VARCHAR(36)   NOT NULL PRIMARY KEY,

  entity_type   VARCHAR(50)   NOT NULL, -- 'vessel_application' | 'organization'
  entity_id     VARCHAR(36)   NOT NULL,

  action        VARCHAR(50)   NOT NULL, -- 'submitted' | 'approved' | 'rejected' | 'withdrawn' | 'reopened'
  old_status    VARCHAR(50)   NULL,
  new_status    VARCHAR(50)   NULL,

  performed_by  VARCHAR(36)   NULL,     -- users.id; NULL if system-generated
  notes         VARCHAR(1000) NULL,

  created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_audit_logs_entity (entity_type, entity_id),
  FOREIGN KEY (performed_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;