package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/dyingfirst/subscribe-server/internal/config"
)

func setupApp(t *testing.T) (*fiber.App, *config.Config) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
app:
  listen: ":8080"
  api_key: "testsecret"
servers:
  - id: "srv_1"
    name: "TestServer"
    host: "test.example.com"
    port: 443
    security: "tls"
    transport: "tcp"
users:
  - id: "usr_1"
    name: "alice"
    uuid: "uuid-alice-1234"
    token: "sub_token_alice"
    active: true
  - id: "usr_2"
    name: "bob"
    uuid: "uuid-bob-5678"
    token: "sub_token_bob"
    active: false
  - id: "usr_3"
    name: "expired"
    uuid: "uuid-expired"
    token: "sub_token_expired"
    active: true
    expires_at: "2020-01-01T00:00:00Z"
  - id: "usr_4"
    name: "noservers"
    uuid: "uuid-noservers"
    token: "sub_token_noservers"
    active: true
    server_ids: ["nonexistent"]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	h := New(cfg)

	app.Get("/sub/:token", h.GetSubscription)

	api := app.Group("/api/v1", h.APIKeyMiddleware)
	api.Get("/users", h.ListUsers)
	api.Post("/users", h.CreateUser)
	api.Delete("/users/:id", h.DeleteUser)
	api.Get("/servers", h.ListServers)
	api.Post("/servers", h.CreateServer)
	api.Delete("/servers/:id", h.DeleteServer)

	return app, cfg
}

// --- Subscription endpoint ---

func TestSubscription_ValidToken(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/sub/sub_token_alice", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Error("body should not be empty")
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "text/plain" {
		t.Errorf("content-type = %q, want text/plain", ct)
	}
}

func TestSubscription_UnknownToken(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/sub/nonexistent", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestSubscription_InactiveUser(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/sub/sub_token_bob", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 403 {
		t.Errorf("status = %d, want 403", resp.StatusCode)
	}
}

func TestSubscription_ExpiredUser(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/sub/sub_token_expired", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 403 {
		t.Errorf("status = %d, want 403", resp.StatusCode)
	}
}

func TestSubscription_NoServersAvailable(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/sub/sub_token_noservers", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

// --- API Key Middleware ---

func TestAPIKey_ValidBearer(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer testsecret")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestAPIKey_InvalidBearer(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestAPIKey_QueryParam(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/api/v1/users?api_key=testsecret", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestAPIKey_NoKey(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

// --- API CRUD ---

func TestAPI_ListUsers(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer testsecret")
	resp, _ := app.Test(req)

	var users []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&users)
	if len(users) != 4 {
		t.Errorf("users = %d, want 4", len(users))
	}
}

func TestAPI_CreateAndDeleteUser(t *testing.T) {
	app, _ := setupApp(t)

	body := bytes.NewBufferString(`{"name":"newuser"}`)
	req := httptest.NewRequest("POST", "/api/v1/users", body)
	req.Header.Set("Authorization", "Bearer testsecret")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, body = %s", resp.StatusCode, body)
	}

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	if created["name"] != "newuser" {
		t.Errorf("name = %v, want newuser", created["name"])
	}

	id, ok := created["id"].(string)
	if !ok || id == "" {
		t.Fatal("created user should have an id")
	}

	delReq := httptest.NewRequest("DELETE", "/api/v1/users/"+id, nil)
	delReq.Header.Set("Authorization", "Bearer testsecret")
	delResp, _ := app.Test(delReq)
	if delResp.StatusCode != 204 {
		t.Errorf("delete status = %d, want 204", delResp.StatusCode)
	}
}

func TestAPI_ListServers(t *testing.T) {
	app, _ := setupApp(t)

	req := httptest.NewRequest("GET", "/api/v1/servers", nil)
	req.Header.Set("Authorization", "Bearer testsecret")
	resp, _ := app.Test(req)

	var servers []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&servers)
	if len(servers) != 1 {
		t.Errorf("servers = %d, want 1", len(servers))
	}
}

func TestAPI_CreateAndDeleteServer(t *testing.T) {
	app, _ := setupApp(t)

	body := bytes.NewBufferString(`{"name":"New","host":"new.example.com","port":443,"security":"tls","transport":"tcp"}`)
	req := httptest.NewRequest("POST", "/api/v1/servers", body)
	req.Header.Set("Authorization", "Bearer testsecret")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, body = %s", resp.StatusCode, b)
	}

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	id, ok := created["id"].(string)
	if !ok || id == "" {
		t.Fatal("created server should have an id")
	}

	delReq := httptest.NewRequest("DELETE", "/api/v1/servers/"+id, nil)
	delReq.Header.Set("Authorization", "Bearer testsecret")
	delResp, _ := app.Test(delReq)
	if delResp.StatusCode != 204 {
		t.Errorf("delete status = %d, want 204", delResp.StatusCode)
	}
}
