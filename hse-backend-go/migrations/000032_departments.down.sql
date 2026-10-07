ALTER TABLE users
  DROP FOREIGN KEY fk_user_department,
  DROP COLUMN department_id;
DROP TABLE IF EXISTS departments;