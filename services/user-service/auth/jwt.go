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

// AuthClaims defines the structure of the JWT claims for authentication.
type AuthClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role,omitempty"`
	Type   string `json:"type"` // e.g., "full_auth", "2fa_temp"
	jwt.RegisteredClaims
}

// GenerateToken generates a new, full-access JWT for a given user ID and role.
func GenerateToken(userID int64, role string) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &AuthClaims{
		UserID: userID,
		Role:   role,
		Type:   "full_auth",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken parses and validates a JWT string, returning the claims if valid.
func ValidateToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Since we're using RSA, we need to provide the public key for verification.
		return &signKey.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}

// Generate2FATempToken generates a short-lived temporary token for 2FA verification.
func Generate2FATempToken(userID int64) (string, error) {
	// Token expires in 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &AuthClaims{
		UserID: userID,
		Type:   "2fa_temp",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
