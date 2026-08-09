INSERT INTO menus (parent_id, name, path, icon, order_index, requires_org_approval) VALUES
  (NULL, 'Organization Review', '/organization/review', 'clipboardList', 9, FALSE);

INSERT INTO role_menu_access (role_id, menu_id, can_view)
SELECT r.id, m.id, TRUE
FROM roles r, menus m
WHERE m.name = 'Organization Review'
  AND r.name IN ('Admin', 'ANP HSE');
