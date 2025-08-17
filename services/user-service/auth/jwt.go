package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	signKey *rsa.PrivateKey
)

// LoadPrivateKey loads the RSA private key from a file.
func LoadPrivateKey(path string) error {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading private key: %w", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return fmt.Errorf("error parsing private key: %w", err)
	}
	return nil
}

// Claims defines the structure of the JWT claims.
type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT for a given user ID and role.
func GenerateToken(userID int64, role string) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	}

	// Create the token with the RS256 signing algorithm and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the private key
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
