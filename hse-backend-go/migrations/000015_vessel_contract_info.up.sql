-- hse-backend-go/migrations/000015_vessel_contract_info.up.sql
ALTER TABLE vessel_applications
  ADD COLUMN contract_info VARCHAR(500) NULL;