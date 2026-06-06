package vless

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func TestBuildSubscription_MultipleServers(t *testing.T) {
	servers := []model.Server{
		{Host: "s1.example.com", Port: 443, Security: "tls", Transport: "tcp", Name: "S1"},
		{Host: "s2.example.com", Port: 443, Security: "reality", Transport: "ws", Name: "S2"},
	}
	got := BuildSubscription(servers, "user-uuid")

	decoded, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}

	lines := strings.Split(string(decoded), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2", len(lines))
	}

	for i, line := range lines {
		if !strings.HasPrefix(line, "vless://") {
			t.Errorf("line %d: not a vless URI", i)
		}
	}

	if !strings.Contains(lines[0], "s1.example.com") {
		t.Error("first line should contain s1.example.com")
	}
	if !strings.Contains(lines[1], "s2.example.com") {
		t.Error("second line should contain s2.example.com")
	}
}

func TestBuildSubscription_SingleServer(t *testing.T) {
	servers := []model.Server{
		{Host: "only.example.com", Port: 443, Security: "tls", Transport: "tcp", Name: "Only"},
	}
	got := BuildSubscription(servers, "uuid")

	decoded, _ := base64.StdEncoding.DecodeString(got)
	if !strings.Contains(string(decoded), "only.example.com") {
		t.Error("should contain server host")
	}
	if strings.Contains(string(decoded), "\n") {
		t.Error("single server should not have newlines")
	}
}
