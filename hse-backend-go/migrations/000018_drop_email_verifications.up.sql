-- email_verifications sudah tidak dipakai lagi - verifikasi email
-- sekarang disimpan di Redis (key "email_verify:<token>", TTL 24h)
-- lihat handlers/authHandler.go (SignUp, VerifyEmail)
DROP TABLE IF EXISTS email_verifications;