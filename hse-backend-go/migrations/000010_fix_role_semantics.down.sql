INSERT INTO role_menu_access (role_id, menu_id, can_view)
SELECT r.id, m.id, TRUE
FROM roles r, menus m
WHERE r.name = 'Staff'
  AND m.name IN ('My Organization', 'Vessel Entry Application');

DELETE FROM roles WHERE name = 'Unassigned';
