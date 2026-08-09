DELETE FROM role_menu_access WHERE menu_id IN (SELECT id FROM menus WHERE name = 'Organization Review');
DELETE FROM menus WHERE name = 'Organization Review';
