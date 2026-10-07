ALTER TABLE vessel_applications
  DROP COLUMN contract_info;

ALTER TABLE vessel_applications
  ADD COLUMN notification_emails JSON NULL;

ALTER TABLE vessel_applications
  ADD COLUMN contract_status        VARCHAR(30)  NULL,
  ADD COLUMN contract_type          VARCHAR(50)  NULL,
  ADD COLUMN contract_type_other    VARCHAR(255) NULL,
  ADD COLUMN contract_number        VARCHAR(100) NULL,
  ADD COLUMN contract_contact_email VARCHAR(255) NULL;

ALTER TABLE vessel_applications
  ADD COLUMN entry_application_types JSON NULL;