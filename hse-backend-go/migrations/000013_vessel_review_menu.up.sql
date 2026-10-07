-- hse-backend-go/migrations/000013_vessel_review_menu.up.sql
INSERT INTO menus (parent_id, name, path, icon, order_index, requires_org_approval) VALUES
  (NULL, 'Vessel Application Review', '/vessel/review', 'clipboardList', 11, FALSE);

INSERT INTO role_menu_access (role_id, menu_id, can_view)
SELECT r.id, m.id, TRUE
FROM roles r, menus m
WHERE m.name = 'Vessel Application Review'
  AND r.name IN ('Admin', 'ANP HSE');