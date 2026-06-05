package vless

import (
	"encoding/base64"
	"strings"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func BuildSubscription(servers []model.Server, userUUID string) string {
	var lines []string
	for _, s := range servers {
		lines = append(lines, BuildURI(s, userUUID))
	}
	joined := strings.Join(lines, "\n")
	return base64.StdEncoding.EncodeToString([]byte(joined))
}
