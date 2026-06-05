package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func (h *Handler) ListServers(c *fiber.Ctx) error {
	return c.JSON(h.cfg.AllServers())
}

func (h *Handler) CreateServer(c *fiber.Ctx) error {
	var s model.Server
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid json")
	}

	s.ID = fmt.Sprintf("srv_%d", time.Now().UnixMilli())

	if err := s.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	h.cfg.AddServer(s)
	if err := h.cfg.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("save failed")
	}

	return c.Status(fiber.StatusCreated).JSON(s)
}

func (h *Handler) DeleteServer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).SendString("id required")
	}

	if !h.cfg.RemoveServer(id) {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}

	if err := h.cfg.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("save failed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
