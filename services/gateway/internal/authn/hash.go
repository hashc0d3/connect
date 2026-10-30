package authn

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashKey(key string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(key)))
	return hex.EncodeToString(sum[:])
}
