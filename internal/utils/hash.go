package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}
