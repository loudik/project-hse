ALTER TABLE vessel_applications
  ADD COLUMN substituted_from_application_id VARCHAR(36) NULL,
  ADD CONSTRAINT fk_vessel_substituted_from
    FOREIGN KEY (substituted_from_application_id) REFERENCES vessel_applications(id) ON DELETE SET NULL;