package service

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashRefreshToken возвращает SHA-256 hex-строку (64 символа).
// подходит для поиска по колонке refresh_token_hash
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
