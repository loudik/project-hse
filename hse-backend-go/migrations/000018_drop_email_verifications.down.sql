-- rollback: kembalikan struktur tabel seperti semula (migration 000007)
-- catatan: data lama (token yang pernah ada) TIDAK bisa dikembalikan,
-- karena sudah tidak disimpan lagi sejak pindah ke Redis
CREATE TABLE IF NOT EXISTS email_verifications (
  id         VARCHAR(36)  NOT NULL PRIMARY KEY,
  user_id    VARCHAR(36)  NOT NULL,
  token      VARCHAR(100) NOT NULL UNIQUE,
  expires_at DATETIME     NOT NULL,
  used_at    DATETIME     NULL,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;