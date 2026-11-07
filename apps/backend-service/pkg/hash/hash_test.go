package hash

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "SecurePassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hashed == password {
		t.Errorf("hash should not match plain password")
	}

	if !CheckPassword(password, hashed) {
		t.Errorf("expected CheckPassword to return true for correct password")
	}

	if CheckPassword("WrongPassword!", hashed) {
		t.Errorf("expected CheckPassword to return false for incorrect password")
	}
}

