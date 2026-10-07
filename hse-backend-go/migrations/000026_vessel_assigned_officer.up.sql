ALTER TABLE vessel_applications
  ADD COLUMN assigned_hse_officer_id VARCHAR(36) NULL,
  ADD CONSTRAINT fk_vessel_assigned_officer
    FOREIGN KEY (assigned_hse_officer_id) REFERENCES users(id) ON DELETE SET NULL;