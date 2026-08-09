-- ============ ROLES ============
CREATE TABLE IF NOT EXISTS roles (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(50)  NOT NULL UNIQUE,
  description VARCHAR(255) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============ USERS ============
CREATE TABLE IF NOT EXISTS users (
  id            VARCHAR(36)  NOT NULL PRIMARY KEY,
  name          VARCHAR(255) NOT NULL,
  email         VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role_id       INT          NOT NULL,
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (role_id) REFERENCES roles(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============ MENUS (tree, self-referencing parent_id) ============
CREATE TABLE IF NOT EXISTS menus (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  parent_id   INT          NULL,
  name        VARCHAR(100) NOT NULL,
  path        VARCHAR(255) NULL,       -- frontend route, NULL if this is just a folder/group
  icon        VARCHAR(50)  NULL,       -- icon key resolved by the frontend (see navigation-icon.config.jsx)
  order_index INT          NOT NULL DEFAULT 0,
  FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============ ROLE <-> MENU ACCESS ============
CREATE TABLE IF NOT EXISTS role_menu_access (
  role_id   INT     NOT NULL,
  menu_id   INT     NOT NULL,
  can_view  BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (role_id, menu_id),
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============ SEED DATA ============

INSERT INTO roles (name, description) VALUES
  ('Admin',       'Full access to the entire system'),
  ('HSE Officer', 'Reviews and approves/rejects HSE submissions'),
  ('Staff',       'Submits HSE certifications/work permits');

-- Default admin user -> email: admin@hse.local / password: admin123
-- (CHANGE THIS PASSWORD after first login in production)
INSERT INTO users (id, name, email, password_hash, role_id) VALUES
  (UUID(), 'Administrator', 'admin@hse.local',
   '$2b$12$9/CMcSJ7BqjsQ8qxrkBf5uINZi4msdWwNAEqu1sZWEkXRGEDRC6Oq',
   (SELECT id FROM roles WHERE name = 'Admin'));

-- Menu tree
INSERT INTO menus (id, parent_id, name, path, icon, order_index) VALUES
  (1, NULL, 'Dashboard',               '/dashboard',      'dashboard',     1),
  (2, NULL, 'HSE Approval',            NULL,               'shieldCheck',   2),
  (3, 2,    'New Submission',          '/hse/new',         'filePlus',      1),
  (4, 2,    'Submissions List',        '/hse/submissions', 'clipboardList', 2),
  (5, NULL, 'Administration',          NULL,               'settings',      3),
  (6, 5,    'User Management',         '/admin/users',     'users',         1),
  (7, 5,    'Role & Menu Management',  '/admin/access',    'keyRound',      2);

-- Default access per role
-- Admin: access to all menus
INSERT INTO role_menu_access (role_id, menu_id, can_view)
SELECT (SELECT id FROM roles WHERE name = 'Admin'), id, TRUE FROM menus;

-- HSE Officer: Dashboard + Submissions List (for approve/reject), no Administration
INSERT INTO role_menu_access (role_id, menu_id, can_view) VALUES
  ((SELECT id FROM roles WHERE name = 'HSE Officer'), 1, TRUE),
  ((SELECT id FROM roles WHERE name = 'HSE Officer'), 2, TRUE),
  ((SELECT id FROM roles WHERE name = 'HSE Officer'), 4, TRUE);

-- Staff: Dashboard + New Submission only
INSERT INTO role_menu_access (role_id, menu_id, can_view) VALUES
  ((SELECT id FROM roles WHERE name = 'Staff'), 1, TRUE),
  ((SELECT id FROM roles WHERE name = 'Staff'), 2, TRUE),
  ((SELECT id FROM roles WHERE name = 'Staff'), 3, TRUE);
