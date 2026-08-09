-- New internal (Microsoft) users start with ZERO menu access until an
-- Admin manually assigns their real role (Staff / HSE Officer / ANP HSE /
-- Admin) - no role_menu_access rows are created for this role on purpose.
INSERT INTO roles (name, description) VALUES
  ('Unassigned', 'New Microsoft sign-in awaiting role assignment by an Admin')
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- "Staff" is internal ANP personnel (signs in via Microsoft) - they are
-- NOT external vessel operators, so they should not see the applicant-only
-- menus (My Organization / Vessel Entry Application). Only "Operator"
-- (external, self-registered) needs those.
DELETE rma FROM role_menu_access rma
JOIN roles r ON r.id = rma.role_id
JOIN menus m ON m.id = rma.menu_id
WHERE r.name = 'Staff'
  AND m.name IN ('My Organization', 'Vessel Entry Application');
