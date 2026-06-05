package model

import "fmt"

type Server struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Host        string `yaml:"host" json:"host"`
	Port        int    `yaml:"port" json:"port"`
	Security    string `yaml:"security" json:"security"`
	Transport   string `yaml:"transport" json:"transport"`
	Path        string `yaml:"path,omitempty" json:"path,omitempty"`
	HostHeader  string `yaml:"host_header,omitempty" json:"host_header,omitempty"`
	SNI         string `yaml:"sni,omitempty" json:"sni,omitempty"`
	Fingerprint string `yaml:"fingerprint,omitempty" json:"fingerprint,omitempty"`
	Flow        string `yaml:"flow,omitempty" json:"flow,omitempty"`
	PBK         string `yaml:"pbk,omitempty" json:"pbk,omitempty"`
	SID         string `yaml:"sid,omitempty" json:"sid,omitempty"`
	ServiceName string `yaml:"service_name,omitempty" json:"service_name,omitempty"`
}

func (s *Server) Validate() error {
	if s.Host == "" {
		return fmt.Errorf("host is required")
	}
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	switch s.Security {
	case "none", "tls", "reality":
	default:
		return fmt.Errorf("security must be none, tls, or reality")
	}
	switch s.Transport {
	case "tcp", "ws", "grpc", "httpupgrade":
	default:
		return fmt.Errorf("transport must be tcp, ws, grpc, or httpupgrade")
	}
	return nil
}
