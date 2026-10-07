ALTER TABLE vessel_applications
  DROP COLUMN notification_emails,
  DROP COLUMN contract_status,
  DROP COLUMN contract_type,
  DROP COLUMN contract_type_other,
  DROP COLUMN contract_number,
  DROP COLUMN contract_contact_email,
  DROP COLUMN entry_application_types;

ALTER TABLE vessel_applications
  ADD COLUMN contract_info VARCHAR(500) NULL;