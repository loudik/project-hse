CREATE TABLE IF NOT EXISTS vessel_extension_requests (
  id                      VARCHAR(36)  NOT NULL PRIMARY KEY,
  application_id          VARCHAR(36)  NOT NULL,
  current_entry_date_to   DATE         NOT NULL,
  requested_entry_date_to DATE         NOT NULL,
  reason                  TEXT         NOT NULL,
  status                  VARCHAR(20)  NOT NULL DEFAULT 'Pending', -- Pending | Approved | Rejected
  rejection_reason        TEXT         NULL,
  requested_by            VARCHAR(36)  NOT NULL,
  decided_by              VARCHAR(36)  NULL,
  decided_at              DATETIME     NULL,
  created_at              DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (application_id) REFERENCES vessel_applications(id) ON DELETE CASCADE,
  FOREIGN KEY (requested_by) REFERENCES users(id) ON DELETE RESTRICT,
  FOREIGN KEY (decided_by) REFERENCES users(id) ON DELETE SET NULL,
  INDEX idx_extension_application (application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tracks whether the H-10 "authorisation expiring soon" reminder has
-- already been sent for this application's CURRENT entry_date_to, so it
-- doesn't get re-sent every day. Resets to NULL whenever an extension is
-- approved (new date = fresh reminder cycle).
ALTER TABLE vessel_applications
  ADD COLUMN authorisation_reminder_sent_at DATETIME NULL;