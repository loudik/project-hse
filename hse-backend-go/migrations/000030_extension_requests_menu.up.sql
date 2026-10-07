INSERT INTO menus (id, parent_id, name, path, icon, order_index) VALUES
  (13, NULL, 'Extension Requests', '/vessel/extension-requests', 'clipboardList', 6);

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view)
SELECT r.id, 13, TRUE FROM roles r WHERE r.name IN ('Admin', 'ANP HSE');