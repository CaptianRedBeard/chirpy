package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword"

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the hashed password is not empty
	if hashedPassword == "" {
		t.Fatal("expected hashed password to be non-empty")
	}

	// Verify the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("expected password to match, got error: %v", err)
	}
}

func TestHashPassword_Error(t *testing.T) {
	// Test with an empty password
	_, err := HashPassword("")
	if err == nil {
		t.Fatal("expected an error for empty password, got none")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecurePassword"

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test with the correct password
	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Fatalf("expected password to match, got error: %v", err)
	}

	// Test with an incorrect password
	incorrectPassword := "wrongPassword"
	err = CheckPasswordHash(hashedPassword, incorrectPassword)
	if err == nil {
		t.Fatal("expected an error for incorrect password, got none")
	}
}
