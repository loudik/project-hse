ALTER TABLE vessel_applications
  DROP FOREIGN KEY fk_vessel_substituted_from,
  DROP COLUMN substituted_from_application_id;