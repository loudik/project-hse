-- ============ NEW ROLES ============
INSERT INTO roles (name, description) VALUES
  ('ANP HSE', 'Reviews and approves organization profiles and vessel entry applications'),
  ('Operator', 'Vessel/organization operator submitting applications')
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- ============ ORGANIZATIONS (Epic 2, US 2.1) ============
CREATE TABLE IF NOT EXISTS organizations (
  id                  VARCHAR(36)   NOT NULL PRIMARY KEY,
  name                VARCHAR(255)  NOT NULL,
  registration_number VARCHAR(100)  NOT NULL UNIQUE,
  type                VARCHAR(100)  NOT NULL,
  address             VARCHAR(500)  NULL,
  country             VARCHAR(100)  NULL,
  phone_number        VARCHAR(50)   NULL,
  email               VARCHAR(255)  NULL,
  website             VARCHAR(255)  NULL,

  status              VARCHAR(20)   NOT NULL DEFAULT 'Pending', -- Pending | Approved | Rejected

  created_by          VARCHAR(36)   NOT NULL,
  approved_by         VARCHAR(36)   NULL,
  approved_at         DATETIME      NULL,
  rejection_reason    VARCHAR(500)  NULL,

  created_at          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (created_by)  REFERENCES users(id),
  FOREIGN KEY (approved_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============ LINK users -> organizations ============
ALTER TABLE users
  ADD COLUMN organization_id VARCHAR(36) NULL AFTER role_id,
  ADD FOREIGN KEY (organization_id) REFERENCES organizations(id);

-- ============ MENU GATING BY ORGANIZATION APPROVAL ============
ALTER TABLE menus
  ADD COLUMN requires_org_approval BOOLEAN NOT NULL DEFAULT FALSE;

-- Placeholder "Vessel" menu (Epic 3) - page not built yet, but menu entry
-- and gating logic are ready. Visible to Staff/Operator once their
-- organization is approved.
INSERT INTO menus (parent_id, name, path, icon, order_index, requires_org_approval) VALUES
  (NULL, 'Vessel Entry Application', '/vessel', 'shieldCheck', 10, TRUE);

INSERT INTO role_menu_access (role_id, menu_id, can_view)
SELECT r.id, m.id, TRUE
FROM roles r, menus m
WHERE m.name = 'Vessel Entry Application'
  AND r.name IN ('Admin', 'Staff', 'Operator', 'ANP HSE');
