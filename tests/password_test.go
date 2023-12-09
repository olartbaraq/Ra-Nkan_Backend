package all_test

import (
	"testing"

	"github.com/olartbaraq/spectrumshelf/utils"
)

func TestGenerateHashPassword(t *testing.T) {
	password := "secret123"

	hashedPassword, err := utils.GenerateHashPassword(password)
	if err != nil {
		t.Fatalf("Error generating hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Fatal("Generated hash password is empty")
	}

	err = utils.VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("Error verifying password: %v", err)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "secret123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := utils.GenerateHashPassword(password)
	if err != nil {
		t.Fatalf("Error generating hash password: %v", err)
	}

	err = utils.VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("Error verifying correct password: %v", err)
	}

	err = utils.VerifyPassword(wrongPassword, hashedPassword)
	if err == nil {
		t.Fatal("Verification should fail for incorrect password")
	}
}

func TestGenerateHashPasswordError(t *testing.T) {
	// Simulate an error in bcrypt.GenerateFromPassword by passing an empty password
	_, err := utils.GenerateHashPassword("")
	if err == nil {
		t.Fatal("Expected error for empty password, but got nil")
	}
}
