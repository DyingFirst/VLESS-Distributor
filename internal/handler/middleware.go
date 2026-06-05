package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) APIKeyMiddleware(c *fiber.Ctx) error {
	key := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if key == "" {
		key = c.Query("api_key")
	}
	if key != h.cfg.App.APIKey {
		return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
	}
	return c.Next()
}

func (h *Handler) CookieAuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("admin_token")
	if token != h.cfg.App.APIKey {
		return c.Redirect("/admin/login")
	}
	return c.Next()
}

func (h *Handler) AdminLoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func (h *Handler) AdminLogin(c *fiber.Ctx) error {
	password := c.FormValue("password")
	if password != h.cfg.App.APIKey {
		return c.Render("login", fiber.Map{"Error": "Invalid password"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    password,
		MaxAge:   86400,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return c.Redirect("/admin/")
}

func (h *Handler) AdminLogout(c *fiber.Ctx) error {
	c.ClearCookie("admin_token")
	return c.Redirect("/admin/login")
}
