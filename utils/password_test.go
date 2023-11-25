package utils

import (
	"testing"
)

func TestGenerateHashPassword(t *testing.T) {
	password := "secret123"

	hashedPassword, err := GenerateHashPassword(password)
	if err != nil {
		t.Fatalf("Error generating hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Fatal("Generated hash password is empty")
	}

	err = VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("Error verifying password: %v", err)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "secret123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := GenerateHashPassword(password)
	if err != nil {
		t.Fatalf("Error generating hash password: %v", err)
	}

	err = VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("Error verifying correct password: %v", err)
	}

	err = VerifyPassword(wrongPassword, hashedPassword)
	if err == nil {
		t.Fatal("Verification should fail for incorrect password")
	}
}

func TestGenerateHashPasswordError(t *testing.T) {
	// Simulate an error in bcrypt.GenerateFromPassword by passing an empty password
	_, err := GenerateHashPassword("")
	if err == nil {
		t.Fatal("Expected error for empty password, but got nil")
	}
}
