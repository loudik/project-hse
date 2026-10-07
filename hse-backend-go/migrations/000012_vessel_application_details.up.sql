-- ============ Perluas vessel_applications ============
ALTER TABLE vessel_applications
  -- Vessel-2: kondisi & entry info
  ADD COLUMN entry_condition      VARCHAR(30)   NULL,  -- Initial | Extension | Substitute
  ADD COLUMN purpose              TEXT          NULL,
  ADD COLUMN entry_date_from      DATE          NULL,
  ADD COLUMN entry_date_to        DATE          NULL,
  ADD COLUMN entry_type           VARCHAR(20)   NULL,  -- single | multiple

  -- Vessel 3: scope & deskripsi vessel
  ADD COLUMN scope_of_work        JSON          NULL,  -- array of selected vessel type options
  ADD COLUMN scope_of_work_other  VARCHAR(255)  NULL,
  ADD COLUMN operation_mode       VARCHAR(50)   NULL,  -- DP1 | DP2 | DP3 | Non-DP
  ADD COLUMN vessel_imo_number    VARCHAR(50)   NULL,
  ADD COLUMN vessel_owner         VARCHAR(255)  NULL,
  ADD COLUMN flag_state           VARCHAR(100)  NULL,
  ADD COLUMN port_of_registry     VARCHAR(150)  NULL,
  ADD COLUMN classification_society VARCHAR(150) NULL,
  ADD COLUMN class_id_number      VARCHAR(100)  NULL,
  ADD COLUMN length_overall       VARCHAR(50)   NULL,
  ADD COLUMN draft_value          VARCHAR(50)   NULL,  -- "draft" is a reserved-ish word, use draft_value
  ADD COLUMN gross_tonnage        VARCHAR(50)   NULL,
  ADD COLUMN call_sign            VARCHAR(50)   NULL,

  -- Vessel-5: declarations
  ADD COLUMN no_major_deficiencies      BOOLEAN NULL,
  ADD COLUMN no_detention_12_months     BOOLEAN NULL,
  ADD COLUMN safety_equipment_operational BOOLEAN NULL,
  ADD COLUMN firefighting_operational   BOOLEAN NULL,
  ADD COLUMN lifesaving_operational     BOOLEAN NULL,
  ADD COLUMN crew_count               INT     NULL,
  ADD COLUMN survey_crew_count        INT     NULL,
  ADD COLUMN cargo_onboard            BOOLEAN NULL,
  ADD COLUMN cargo_description        TEXT    NULL,
  ADD COLUMN hazardous_cargo          BOOLEAN NULL,
  ADD COLUMN waste_discharge          BOOLEAN NULL,
  ADD COLUMN oily_waste_onboard       BOOLEAN NULL,
  ADD COLUMN sewage_disposal_required BOOLEAN NULL;

-- ============ Dokumen/file generik (Vessel 1, 4, 5, cover letter, dll) ============
CREATE TABLE IF NOT EXISTS vessel_application_documents (
  id                VARCHAR(36)  NOT NULL PRIMARY KEY,
  application_id    VARCHAR(36)  NOT NULL,

  category          VARCHAR(50)  NOT NULL, -- regulatory_document | statutory_certificate | supporting_document | cover_letter | crew_document | cargo_document | clearance_document
  document_key      VARCHAR(100) NOT NULL, -- e.g. "vessel_safety_case", "registration_certificate", "immigration_clearance"
  label             VARCHAR(255) NULL,     -- human-readable name, useful for "Other, please specify"

  not_applicable    BOOLEAN      NOT NULL DEFAULT FALSE,
  na_reason         VARCHAR(500) NULL,

  date_issued       DATE         NULL,
  date_expired      DATE         NULL,

  file_path         VARCHAR(500) NULL,
  file_name         VARCHAR(255) NULL,
  file_size         INT          NULL,

  uploaded_at       DATETIME     NULL,
  created_at        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (application_id) REFERENCES vessel_applications(id) ON DELETE CASCADE,
  UNIQUE KEY uniq_app_doc (application_id, category, document_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;