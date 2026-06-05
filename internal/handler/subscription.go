package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/dyingfirst/subscribe-server/internal/vless"
)

func (h *Handler) GetSubscription(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).SendString("token required")
	}

	user := h.cfg.FindUserByToken(token)
	if user == nil {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}

	if !user.Active || user.IsExpired() {
		return c.Status(fiber.StatusForbidden).SendString("subscription expired or inactive")
	}

	servers := h.cfg.GetServersForUser(user)
	if len(servers) == 0 {
		return c.Status(fiber.StatusNotFound).SendString("no servers available")
	}

	result := vless.BuildSubscription(servers, user.UUID)

	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "attachment; filename=subscription.txt")
	return c.SendString(result)
}
