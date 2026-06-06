package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

const testYAML = `
app:
  listen: ":9090"
  api_key: "testkey"
servers:
  - id: "srv_1"
    name: "Test"
    host: "test.example.com"
    port: 443
    security: "tls"
    transport: "tcp"
users:
  - id: "usr_1"
    name: "alice"
    uuid: "uuid-alice"
    token: "sub_token_alice"
    active: true
  - id: "usr_2"
    name: "bob"
    uuid: "uuid-bob"
    token: "sub_token_bob"
    active: false
    server_ids:
      - "srv_1"
`

func TestLoad_Valid(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.App.Listen != ":9090" {
		t.Errorf("listen = %q, want :9090", cfg.App.Listen)
	}
	if cfg.App.APIKey != "testkey" {
		t.Errorf("apikey = %q, want testkey", cfg.App.APIKey)
	}
	if cfg.configPath != path {
		t.Errorf("configPath = %q, want %q", cfg.configPath, path)
	}
	if len(cfg.Users) != 2 {
		t.Errorf("users = %d, want 2", len(cfg.Users))
	}
	if len(cfg.Servers) != 1 {
		t.Errorf("servers = %d, want 1", len(cfg.Servers))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("missing file should return error")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTestConfig(t, "{{invalid yaml")
	_, err := Load(path)
	if err == nil {
		t.Error("invalid YAML should return error")
	}
}

func TestLoad_DefaultListen(t *testing.T) {
	path := writeTestConfig(t, `app: {api_key: "k"}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.App.Listen != ":8080" {
		t.Errorf("default listen = %q, want :8080", cfg.App.Listen)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	cfg.AddUser(model.User{ID: "usr_3", Name: "charlie", UUID: "uuid-3", Token: "sub_3", Active: true})
	cfg.AddServer(model.Server{ID: "srv_2", Name: "S2", Host: "s2.example.com", Port: 443, Security: "tls", Transport: "tcp"})

	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Users) != 3 {
		t.Errorf("users after save = %d, want 3", len(loaded.Users))
	}
	if len(loaded.Servers) != 2 {
		t.Errorf("servers after save = %d, want 2", len(loaded.Servers))
	}
}

func TestFindUserByToken(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	u := cfg.FindUserByToken("sub_token_alice")
	if u == nil {
		t.Fatal("should find alice")
	}
	if u.Name != "alice" {
		t.Errorf("name = %q, want alice", u.Name)
	}
	if cfg.FindUserByToken("nonexistent") != nil {
		t.Error("should not find unknown token")
	}
}

func TestFindUserByID(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	u := cfg.FindUserByID("usr_1")
	if u == nil {
		t.Fatal("should find usr_1")
	}
	if u.Name != "alice" {
		t.Errorf("name = %q, want alice", u.Name)
	}
	if cfg.FindUserByID("nonexistent") != nil {
		t.Error("should not find unknown id")
	}
}

func TestRemoveUser(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	if !cfg.RemoveUser("usr_1") {
		t.Error("should remove existing user")
	}
	if len(cfg.AllUsers()) != 1 {
		t.Errorf("users = %d, want 1", len(cfg.AllUsers()))
	}
	if cfg.RemoveUser("nonexistent") {
		t.Error("should not remove unknown user")
	}
}

func TestFindServerByID(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	s := cfg.FindServerByID("srv_1")
	if s == nil {
		t.Fatal("should find srv_1")
	}
	if s.Name != "Test" {
		t.Errorf("name = %q, want Test", s.Name)
	}
	if cfg.FindServerByID("nonexistent") != nil {
		t.Error("should not find unknown id")
	}
}

func TestRemoveServer(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	if !cfg.RemoveServer("srv_1") {
		t.Error("should remove existing server")
	}
	if len(cfg.AllServers()) != 0 {
		t.Errorf("servers = %d, want 0", len(cfg.AllServers()))
	}
	if cfg.RemoveServer("nonexistent") {
		t.Error("should not remove unknown server")
	}
}

func TestGetServersForUser_All(t *testing.T) {
	path := writeTestConfig(t, `
app: {api_key: "k"}
servers:
  - {id: "s1", name: "A", host: "a", port: 443, security: "tls", transport: "tcp"}
  - {id: "s2", name: "B", host: "b", port: 443, security: "tls", transport: "tcp"}
users:
  - {id: "u1", name: "alice", uuid: "x", token: "t1", active: true}
`)
	cfg, _ := Load(path)
	u := cfg.FindUserByID("u1")

	servers := cfg.GetServersForUser(u)
	if len(servers) != 2 {
		t.Errorf("servers = %d, want 2 (all)", len(servers))
	}
}

func TestGetServersForUser_Filtered(t *testing.T) {
	path := writeTestConfig(t, `
app: {api_key: "k"}
servers:
  - {id: "s1", name: "A", host: "a", port: 443, security: "tls", transport: "tcp"}
  - {id: "s2", name: "B", host: "b", port: 443, security: "tls", transport: "tcp"}
users:
  - {id: "u1", name: "alice", uuid: "x", token: "t1", active: true, server_ids: ["s2"]}
`)
	cfg, _ := Load(path)
	u := cfg.FindUserByID("u1")

	servers := cfg.GetServersForUser(u)
	if len(servers) != 1 {
		t.Fatalf("servers = %d, want 1", len(servers))
	}
	if servers[0].ID != "s2" {
		t.Errorf("server id = %q, want s2", servers[0].ID)
	}
}

func TestGetServersForUser_Empty(t *testing.T) {
	path := writeTestConfig(t, `
app: {api_key: "k"}
servers: []
users:
  - {id: "u1", name: "alice", uuid: "x", token: "t1", active: true, server_ids: ["s1"]}
`)
	cfg, _ := Load(path)
	u := cfg.FindUserByID("u1")

	servers := cfg.GetServersForUser(u)
	if len(servers) != 0 {
		t.Errorf("servers = %d, want 0", len(servers))
	}
}

func TestAllUsers_ReturnsCopy(t *testing.T) {
	path := writeTestConfig(t, testYAML)
	cfg, _ := Load(path)

	users := cfg.AllUsers()
	users[0].Name = "modified"

	original := cfg.FindUserByID("usr_1")
	if original.Name == "modified" {
		t.Error("AllUsers should return a copy")
	}
}

func TestIsExpired(t *testing.T) {
	path := writeTestConfig(t, `
app: {api_key: "k"}
users:
  - {id: "u1", name: "alice", uuid: "x", token: "t1", active: true, expires_at: "2020-01-01T00:00:00Z"}
`)
	cfg, _ := Load(path)
	u := cfg.FindUserByID("u1")
	if !u.IsExpired() {
		t.Error("2020 date should be expired")
	}
}

// ensure time import is used
var _ time.Time
