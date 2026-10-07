ALTER TABLE vessel_applications
  DROP FOREIGN KEY fk_vessel_assigned_officer,
  DROP COLUMN assigned_hse_officer_id;