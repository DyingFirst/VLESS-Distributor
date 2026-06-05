package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

type createUserRequest struct {
	Name      string   `json:"name"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	ServerIDs []string `json:"server_ids,omitempty"`
}

func (h *Handler) ListUsers(c *fiber.Ctx) error {
	return c.JSON(h.cfg.AllUsers())
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid json")
	}

	u := model.User{
		ID:        fmt.Sprintf("usr_%s", model.GenerateToken()[4:16]),
		Name:      req.Name,
		UUID:      model.GenerateUUID(),
		Token:     model.GenerateToken(),
		Active:    true,
		ServerIDs: req.ServerIDs,
	}

	if err := u.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if req.ExpiresAt != "" {
		t, err := parseTime(req.ExpiresAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid expires_at format, use RFC3339")
		}
		u.ExpiresAt = &t
	}

	h.cfg.AddUser(u)
	if err := h.cfg.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("save failed")
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).SendString("id required")
	}

	if !h.cfg.RemoveUser(id) {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}

	if err := h.cfg.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("save failed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
