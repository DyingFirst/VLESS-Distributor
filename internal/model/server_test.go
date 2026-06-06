package model

import "testing"

func TestServerValidate_Valid(t *testing.T) {
	cases := []struct {
		name string
		s    Server
	}{
		{"tls/tcp", Server{Host: "h", Port: 443, Security: "tls", Transport: "tcp"}},
		{"reality/tcp", Server{Host: "h", Port: 443, Security: "reality", Transport: "tcp"}},
		{"none/ws", Server{Host: "h", Port: 80, Security: "none", Transport: "ws"}},
		{"tls/grpc", Server{Host: "h", Port: 443, Security: "tls", Transport: "grpc"}},
		{"tls/httpupgrade", Server{Host: "h", Port: 443, Security: "tls", Transport: "httpupgrade"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.s.Validate(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestServerValidate_MissingHost(t *testing.T) {
	s := Server{Host: "", Port: 443, Security: "tls", Transport: "tcp"}
	if err := s.Validate(); err == nil {
		t.Error("empty host should fail")
	}
}

func TestServerValidate_BadPort(t *testing.T) {
	cases := []struct {
		name string
		port int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too large", 65536},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Server{Host: "h", Port: tc.port, Security: "tls", Transport: "tcp"}
			if err := s.Validate(); err == nil {
				t.Error("bad port should fail")
			}
		})
	}
}

func TestServerValidate_BadSecurity(t *testing.T) {
	s := Server{Host: "h", Port: 443, Security: "bad", Transport: "tcp"}
	if err := s.Validate(); err == nil {
		t.Error("bad security should fail")
	}
}

func TestServerValidate_BadTransport(t *testing.T) {
	s := Server{Host: "h", Port: 443, Security: "tls", Transport: "bad"}
	if err := s.Validate(); err == nil {
		t.Error("bad transport should fail")
	}
}
