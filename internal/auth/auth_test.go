package auth_test

import (
	"testing"
	"time"

	"github.com/brobridge/sentinel/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestVerifyPassword_Correct(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword: %v", err)
	}
	if err := auth.VerifyPassword("secret123", string(hash)); err != nil {
		t.Errorf("VerifyPassword() should pass for correct password: %v", err)
	}
}

func TestVerifyPassword_Wrong(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err := auth.VerifyPassword("wrong", string(hash)); err == nil {
		t.Error("VerifyPassword() should fail for wrong password")
	}
}

func TestHashPassword(t *testing.T) {
	h, err := auth.HashPassword("mypassword")
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	if err := auth.VerifyPassword("mypassword", h); err != nil {
		t.Errorf("HashPassword output should pass VerifyPassword: %v", err)
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	secret := []byte("test-secret-32-bytes-long-enough!")

	token, err := auth.GenerateToken("admin", secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := auth.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken() error: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("username = %q, want admin", claims.Username)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, _ := auth.GenerateToken("admin", []byte("secret-a"), time.Hour)
	if _, err := auth.ValidateToken(token, []byte("secret-b")); err == nil {
		t.Error("ValidateToken() should fail with wrong secret")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	secret := []byte("test-secret-32-bytes-long-enough!")
	token, _ := auth.GenerateToken("admin", secret, -time.Second)
	if _, err := auth.ValidateToken(token, secret); err == nil {
		t.Error("ValidateToken() should fail for expired token")
	}
}
