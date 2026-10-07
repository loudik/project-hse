package utils

import (
	"context"
	"time"

	"hse-backend-go/db"
)

const blacklistPrefix = "blacklist:jti:"

func BlacklistToken(jti string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	ctx := context.Background()
	return db.RDB.Set(ctx, blacklistPrefix+jti, "1", ttl).Err()
}

// IsTokenBlacklisted checks whether a jti was revoked (e.g. via sign-out).
func IsTokenBlacklisted(jti string) bool {
	if jti == "" {
		return false
	}
	ctx := context.Background()
	exists, err := db.RDB.Exists(ctx, blacklistPrefix+jti).Result()
	if err != nil {
		// fail open on Redis errors so an outage doesn't lock out every user
		return false
	}
	return exists > 0
}
