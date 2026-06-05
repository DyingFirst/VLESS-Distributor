package model

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

type User struct {
	ID        string     `yaml:"id" json:"id"`
	Name      string     `yaml:"name" json:"name"`
	UUID      string     `yaml:"uuid" json:"uuid"`
	Token     string     `yaml:"token" json:"token"`
	Active    bool       `yaml:"active" json:"active"`
	ExpiresAt *time.Time `yaml:"expires_at,omitempty" json:"expires_at,omitempty"`
	ServerIDs []string   `yaml:"server_ids,omitempty" json:"server_ids,omitempty"`
}

func (u *User) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func GenerateUUID() string {
	b := make([]byte, 16)
	cryptorand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func GenerateToken() string {
	b := make([]byte, 24)
	cryptorand.Read(b)
	return "sub_" + base64.RawURLEncoding.EncodeToString(b)
}
