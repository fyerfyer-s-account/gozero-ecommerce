package cryptx

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}
