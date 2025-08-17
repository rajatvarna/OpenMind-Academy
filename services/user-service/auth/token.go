package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken generates a random, secure token.
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
