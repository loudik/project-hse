UPDATE roles SET description = 'Akses penuh ke seluruh sistem' WHERE name = 'Admin';
UPDATE roles SET description = 'Meninjau dan menyetujui/menolak pengajuan HSE' WHERE name = 'HSE Officer';
UPDATE roles SET description = 'Mengajukan sertifikasi/izin kerja HSE' WHERE name = 'Staff';

UPDATE menus SET name = 'Ajukan Sertifikasi'    WHERE name = 'New Submission';
UPDATE menus SET name = 'Daftar Pengajuan'      WHERE name = 'Submissions List';
UPDATE menus SET name = 'Administrasi'          WHERE name = 'Administration';
UPDATE menus SET name = 'Manajemen User'        WHERE name = 'User Management';
UPDATE menus SET name = 'Manajemen Role & Menu' WHERE name = 'Role & Menu Management';

UPDATE hse_submissions SET status = 'Diajukan'  WHERE status = 'Submitted';
UPDATE hse_submissions SET status = 'Disetujui' WHERE status = 'Approved';
UPDATE hse_submissions SET status = 'Ditolak'   WHERE status = 'Rejected';

ALTER TABLE hse_submissions ALTER COLUMN status SET DEFAULT 'Diajukan';
