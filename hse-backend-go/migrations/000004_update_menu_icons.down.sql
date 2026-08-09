-- Note: at this point in a rollback sequence, migration 000005's down has
-- already reverted menu names back to Indonesian - match against those.
UPDATE menus SET icon = 'LayoutDashboard' WHERE name = 'Dashboard';
UPDATE menus SET icon = 'ShieldCheck'     WHERE name = 'HSE Approval';
UPDATE menus SET icon = 'FilePlus'        WHERE name = 'Ajukan Sertifikasi';
UPDATE menus SET icon = 'ClipboardList'   WHERE name = 'Daftar Pengajuan';
UPDATE menus SET icon = 'Settings'        WHERE name = 'Administrasi';
UPDATE menus SET icon = 'Users'           WHERE name = 'Manajemen User';
UPDATE menus SET icon = 'KeyRound'        WHERE name = 'Manajemen Role & Menu';
