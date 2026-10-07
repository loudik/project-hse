CREATE TABLE IF NOT EXISTS vessel_authorisation_letters (
  id                 VARCHAR(36)   NOT NULL PRIMARY KEY,
  seq                INT           NOT NULL AUTO_INCREMENT UNIQUE,
  application_id     VARCHAR(36)   NOT NULL,

  reference_no       VARCHAR(50)   NOT NULL UNIQUE,
  facility_name      VARCHAR(255)  NOT NULL,
  decree_law_number  VARCHAR(255)  NOT NULL,
  decree_law_clause  VARCHAR(255)  NOT NULL,
  inspection_date    DATE          NULL,
  psc_name           VARCHAR(255)  NOT NULL,

  signature_path     VARCHAR(500)  NOT NULL,
  digital_hash       VARCHAR(128)  NOT NULL,

  docx_path          VARCHAR(500)  NOT NULL,
  pdf_path            VARCHAR(500)  NULL,

  issued_by          VARCHAR(36)   NOT NULL,
  issued_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (application_id) REFERENCES vessel_applications(id) ON DELETE CASCADE,
  FOREIGN KEY (issued_by) REFERENCES users(id) ON DELETE RESTRICT,
  INDEX idx_authorisation_application (application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;