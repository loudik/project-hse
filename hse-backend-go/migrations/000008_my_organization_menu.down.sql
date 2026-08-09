DELETE FROM role_menu_access WHERE menu_id IN (SELECT id FROM menus WHERE name = 'My Organization');
DELETE FROM menus WHERE name = 'My Organization';
