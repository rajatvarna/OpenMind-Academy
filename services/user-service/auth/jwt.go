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

func init() {
	// In a real app, the path to the key should be configurable.
	keyBytes, err := ioutil.ReadFile("../secrets/jwtRS256.key")
	if err != nil {
		panic(fmt.Sprintf("Error reading private key: %s", err))
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		panic(fmt.Sprintf("Error parsing private key: %s", err))
	}
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
