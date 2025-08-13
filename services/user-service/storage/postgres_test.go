package storage

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	password := "my-strong-password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test case 1: Correct password
	if !CheckPassword(string(hashedPassword), password) {
		t.Errorf("CheckPassword failed: expected true for correct password, got false")
	}

	// Test case 2: Incorrect password
	if CheckPassword(string(hashedPassword), "wrong-password") {
		t.Errorf("CheckPassword failed: expected false for incorrect password, got true")
	}
}
