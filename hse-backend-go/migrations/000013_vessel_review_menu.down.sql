-- hse-backend-go/migrations/000013_vessel_review_menu.down.sql
DELETE FROM role_menu_access WHERE menu_id IN (SELECT id FROM menus WHERE name = 'Vessel Application Review');
DELETE FROM menus WHERE name = 'Vessel Application Review';