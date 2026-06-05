package vless

import (
	"fmt"
	"net/url"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func BuildURI(s model.Server, userUUID string) string {
	u := fmt.Sprintf("vless://%s@%s:%d?", userUUID, s.Host, s.Port)

	params := url.Values{}
	params.Set("encryption", "none")
	params.Set("security", s.Security)
	params.Set("type", s.Transport)

	if s.SNI != "" {
		params.Set("sni", s.SNI)
	}
	if s.Flow != "" {
		params.Set("flow", s.Flow)
	}
	if s.Fingerprint != "" {
		params.Set("fp", s.Fingerprint)
	}

	if s.Security == "reality" {
		if s.PBK != "" {
			params.Set("pbk", s.PBK)
		}
		if s.SID != "" {
			params.Set("sid", s.SID)
		}
	}

	switch s.Transport {
	case "ws", "httpupgrade":
		if s.Path != "" {
			params.Set("path", s.Path)
		}
		if s.HostHeader != "" {
			params.Set("host", s.HostHeader)
		}
	case "grpc":
		if s.ServiceName != "" {
			params.Set("serviceName", s.ServiceName)
		}
	}

	fragment := url.PathEscape(s.Name)
	return u + params.Encode() + "#" + fragment
}
