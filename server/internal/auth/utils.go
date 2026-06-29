package auth

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func isValidPassword(password string) bool {
	if password == "" || strings.TrimSpace(password) == "" {
		return false
	}
	if len(password) < 8 {
		return false
	}
	return true
}

func hashToken(tokenstring string) string {
	hash := sha256.Sum256([]byte(tokenstring))
	hashString := fmt.Sprintf("%x", hash)
	return hashString
}
