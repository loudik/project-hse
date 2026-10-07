-- Menu "Dashboard" (id 1) sudah ada dari migration paling awal (000002) -
-- itu cuma placeholder yang mengarah ke /dashboard tanpa halaman nyata,
-- sampai sekarang kita baru bikin halamannya. Jadi di sini kita TIDAK
-- bikin menu baru, cuma pastikan role ANP HSE (selain Admin yang sudah
-- ada) bisa lihat menu itu.
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view)
SELECT r.id, m.id, TRUE
FROM roles r, menus m
WHERE m.name = 'Dashboard'
  AND r.name IN ('Admin', 'ANP HSE');