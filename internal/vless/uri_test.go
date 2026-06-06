package vless

import (
	"net/url"
	"testing"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func TestBuildURI_TLSSimple(t *testing.T) {
	s := model.Server{Host: "us1.example.com", Port: 443, Security: "tls", Transport: "tcp", Name: "US East"}
	got := BuildURI(s, "user-uuid-1234")

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse URI: %v", err)
	}

	if u.Scheme != "vless" {
		t.Errorf("scheme = %q, want vless", u.Scheme)
	}
	if u.User.Username() != "user-uuid-1234" {
		t.Errorf("user = %q, want user-uuid-1234", u.User.Username())
	}
	if u.Hostname() != "us1.example.com" {
		t.Errorf("host = %q, want us1.example.com", u.Hostname())
	}
	if u.Port() != "443" {
		t.Errorf("port = %q, want 443", u.Port())
	}
	if u.Fragment != "US East" {
		t.Errorf("fragment = %q, want US East", u.Fragment)
	}

	q := u.Query()
	if q.Get("security") != "tls" {
		t.Errorf("security = %q, want tls", q.Get("security"))
	}
	if q.Get("type") != "tcp" {
		t.Errorf("type = %q, want tcp", q.Get("type"))
	}
	if q.Get("encryption") != "none" {
		t.Errorf("encryption = %q, want none", q.Get("encryption"))
	}
}

func TestBuildURI_RealityWithFlow(t *testing.T) {
	s := model.Server{
		Host:        "nl1.example.com",
		Port:        443,
		Security:    "reality",
		Transport:   "tcp",
		Flow:        "xtls-rprx-vision",
		SNI:         "www.microsoft.com",
		Fingerprint: "chrome",
		PBK:         "pubkey123",
		SID:         "abc12345",
		Name:        "NL Reality",
	}
	got := BuildURI(s, "uuid-456")

	u, _ := url.Parse(got)
	q := u.Query()

	if q.Get("security") != "reality" {
		t.Errorf("security = %q, want reality", q.Get("security"))
	}
	if q.Get("flow") != "xtls-rprx-vision" {
		t.Errorf("flow = %q, want xtls-rprx-vision", q.Get("flow"))
	}
	if q.Get("sni") != "www.microsoft.com" {
		t.Errorf("sni = %q, want www.microsoft.com", q.Get("sni"))
	}
	if q.Get("fp") != "chrome" {
		t.Errorf("fp = %q, want chrome", q.Get("fp"))
	}
	if q.Get("pbk") != "pubkey123" {
		t.Errorf("pbk = %q, want pubkey123", q.Get("pbk"))
	}
	if q.Get("sid") != "abc12345" {
		t.Errorf("sid = %q, want abc12345", q.Get("sid"))
	}
}

func TestBuildURI_WebSocket(t *testing.T) {
	s := model.Server{
		Host:       "ws1.example.com",
		Port:       443,
		Security:   "tls",
		Transport:  "ws",
		Path:       "/ws",
		HostHeader: "ws1.example.com",
		Name:       "WS Server",
	}
	got := BuildURI(s, "uuid-789")

	u, _ := url.Parse(got)
	q := u.Query()

	if q.Get("type") != "ws" {
		t.Errorf("type = %q, want ws", q.Get("type"))
	}
	if q.Get("path") != "/ws" {
		t.Errorf("path = %q, want /ws", q.Get("path"))
	}
	if q.Get("host") != "ws1.example.com" {
		t.Errorf("host = %q, want ws1.example.com", q.Get("host"))
	}
}

func TestBuildURI_gRPC(t *testing.T) {
	s := model.Server{
		Host:        "grpc1.example.com",
		Port:        443,
		Security:    "tls",
		Transport:   "grpc",
		ServiceName: "grpc-service",
		Name:        "gRPC Server",
	}
	got := BuildURI(s, "uuid-grpc")

	u, _ := url.Parse(got)
	q := u.Query()

	if q.Get("type") != "grpc" {
		t.Errorf("type = %q, want grpc", q.Get("type"))
	}
	if q.Get("serviceName") != "grpc-service" {
		t.Errorf("serviceName = %q, want grpc-service", q.Get("serviceName"))
	}
}

func TestBuildURI_HTTPUpgrade(t *testing.T) {
	s := model.Server{
		Host:       "hu.example.com",
		Port:       443,
		Security:   "tls",
		Transport:  "httpupgrade",
		Path:       "/upgrade",
		HostHeader: "hu.example.com",
		Name:       "HTTPUpgrade",
	}
	got := BuildURI(s, "uuid-hu")

	u, _ := url.Parse(got)
	q := u.Query()

	if q.Get("type") != "httpupgrade" {
		t.Errorf("type = %q, want httpupgrade", q.Get("type"))
	}
	if q.Get("path") != "/upgrade" {
		t.Errorf("path = %q, want /upgrade", q.Get("path"))
	}
	if q.Get("host") != "hu.example.com" {
		t.Errorf("host = %q, want hu.example.com", q.Get("host"))
	}
}

func TestBuildURI_EmptyFields(t *testing.T) {
	s := model.Server{Host: "simple.example.com", Port: 80, Security: "none", Transport: "tcp", Name: "Simple"}
	got := BuildURI(s, "uuid")

	u, _ := url.Parse(got)
	q := u.Query()

	if q.Get("sni") != "" {
		t.Errorf("sni should be empty")
	}
	if q.Get("flow") != "" {
		t.Errorf("flow should be empty")
	}
	if q.Get("fp") != "" {
		t.Errorf("fp should be empty")
	}
	if q.Get("path") != "" {
		t.Errorf("path should be empty")
	}
}

func TestBuildURI_NameEscaped(t *testing.T) {
	s := model.Server{Host: "h", Port: 1, Security: "none", Transport: "tcp", Name: "Server Name/With Special"}
	got := BuildURI(s, "uuid")

	u, _ := url.Parse(got)
	if u.Fragment == "" {
		t.Error("fragment should not be empty")
	}
}
