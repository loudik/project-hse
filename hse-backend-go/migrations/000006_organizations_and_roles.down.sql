DELETE FROM role_menu_access WHERE menu_id IN (SELECT id FROM menus WHERE name = 'Vessel Entry Application');
DELETE FROM menus WHERE name = 'Vessel Entry Application';
ALTER TABLE menus DROP COLUMN requires_org_approval;

-- NOTE: the FK constraint name below (users_ibfk_2) assumes this is the
-- second foreign key added to `users` (after role_id's FK from migration
-- 000002). If this DROP fails, run `SHOW CREATE TABLE users;` first to
-- find the actual constraint name for organization_id and adjust below.
ALTER TABLE users DROP FOREIGN KEY users_ibfk_2;
ALTER TABLE users DROP COLUMN organization_id;

DROP TABLE IF EXISTS organizations;

DELETE FROM roles WHERE name IN ('ANP HSE', 'Operator');
