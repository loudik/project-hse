CREATE TABLE IF NOT EXISTS vessel_applications (
  id                     VARCHAR(36)  NOT NULL PRIMARY KEY,
  application_number     VARCHAR(64)  NOT NULL UNIQUE,
  organization_id        VARCHAR(36)  NOT NULL,

  vessel_name            VARCHAR(255) NULL,
  vessel_type            VARCHAR(100) NULL,
  proposed_activity      TEXT         NULL,
  planned_arrival_date   DATE         NULL,
  planned_departure_date DATE         NULL,

  status                 VARCHAR(20)  NOT NULL DEFAULT 'Draft', -- Draft | Submitted | Approved | Rejected | Withdrawn

  created_by             VARCHAR(36)  NOT NULL,
  submitted_at           DATETIME     NULL,
  created_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  FOREIGN KEY (organization_id) REFERENCES organizations(id),
  FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;