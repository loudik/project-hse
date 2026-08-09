-- Role descriptions
UPDATE roles SET description = 'Full access to the entire system' WHERE name = 'Admin';
UPDATE roles SET description = 'Reviews and approves/rejects HSE submissions' WHERE name = 'HSE Officer';
UPDATE roles SET description = 'Submits HSE certifications/work permits' WHERE name = 'Staff';

-- Menu names
UPDATE menus SET name = 'New Submission'         WHERE name = 'Ajukan Sertifikasi';
UPDATE menus SET name = 'Submissions List'       WHERE name = 'Daftar Pengajuan';
UPDATE menus SET name = 'Administration'         WHERE name = 'Administrasi';
UPDATE menus SET name = 'User Management'        WHERE name = 'Manajemen User';
UPDATE menus SET name = 'Role & Menu Management' WHERE name = 'Manajemen Role & Menu';

-- Existing submission statuses (Indonesian -> English), in case any test
-- data was already created before this change
UPDATE hse_submissions SET status = 'Submitted' WHERE status = 'Diajukan';
UPDATE hse_submissions SET status = 'Approved'  WHERE status = 'Disetujui';
UPDATE hse_submissions SET status = 'Rejected'  WHERE status = 'Ditolak';

-- New default for the status column going forward
ALTER TABLE hse_submissions ALTER COLUMN status SET DEFAULT 'Submitted';
