package model

import (
	"strings"
	"testing"
	"time"
)

func TestIsExpired_NilExpiresAt(t *testing.T) {
	u := User{ExpiresAt: nil}
	if u.IsExpired() {
		t.Error("nil ExpiresAt should not be expired")
	}
}

func TestIsExpired_Future(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	u := User{ExpiresAt: &future}
	if u.IsExpired() {
		t.Error("future ExpiresAt should not be expired")
	}
}

func TestIsExpired_Past(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	u := User{ExpiresAt: &past}
	if !u.IsExpired() {
		t.Error("past ExpiresAt should be expired")
	}
}

func TestValidate_EmptyName(t *testing.T) {
	u := User{Name: ""}
	if err := u.Validate(); err == nil {
		t.Error("empty name should fail validation")
	}
}

func TestValidate_ValidName(t *testing.T) {
	u := User{Name: "alice"}
	if err := u.Validate(); err != nil {
		t.Errorf("valid name should pass: %v", err)
	}
}

func TestGenerateUUID_Format(t *testing.T) {
	uuid := GenerateUUID()
	if len(uuid) != 36 {
		t.Errorf("UUID length = %d, want 36", len(uuid))
	}
	if strings.Count(uuid, "-") != 4 {
		t.Errorf("UUID should have 4 dashes, got %d", strings.Count(uuid, "-"))
	}
}

func TestGenerateUUID_Uniqueness(t *testing.T) {
	a := GenerateUUID()
	b := GenerateUUID()
	if a == b {
		t.Error("two generated UUIDs should differ")
	}
}

func TestGenerateToken_Prefix(t *testing.T) {
	token := GenerateToken()
	if !strings.HasPrefix(token, "sub_") {
		t.Errorf("token should have sub_ prefix, got %q", token)
	}
}

func TestGenerateToken_Uniqueness(t *testing.T) {
	a := GenerateToken()
	b := GenerateToken()
	if a == b {
		t.Error("two generated tokens should differ")
	}
}
