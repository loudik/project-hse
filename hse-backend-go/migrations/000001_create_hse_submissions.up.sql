CREATE TABLE IF NOT EXISTS hse_submissions (
  id                  VARCHAR(36)   NOT NULL PRIMARY KEY,
  cert_number         VARCHAR(64)   NOT NULL UNIQUE,

  applicant_name      VARCHAR(255)  NOT NULL,
  applicant_position  VARCHAR(255)  NULL,
  department          VARCHAR(255)  NOT NULL,
  work_type           VARCHAR(255)  NOT NULL,
  location            VARCHAR(255)  NOT NULL,
  start_date          DATE          NOT NULL,
  end_date            DATE          NOT NULL,
  description         TEXT          NOT NULL,
  hazards             TEXT          NOT NULL,
  controls            TEXT          NOT NULL,

  photos              JSON          NULL,
  applicant_signature LONGTEXT      NOT NULL,

  status              VARCHAR(20)   NOT NULL DEFAULT 'Submitted',

  approver_name       VARCHAR(255)  NULL,
  approver_note       TEXT          NULL,
  approver_signature  LONGTEXT      NULL,

  submitted_at        DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  decided_at          DATETIME      NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
