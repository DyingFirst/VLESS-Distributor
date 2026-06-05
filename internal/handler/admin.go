package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

func (h *Handler) AdminDashboard(c *fiber.Ctx) error {
	users := h.cfg.AllUsers()
	servers := h.cfg.AllServers()
	activeUsers := 0
	for _, u := range users {
		if u.Active && !u.IsExpired() {
			activeUsers++
		}
	}
	return c.Render("dashboard", fiber.Map{
		"UserCount":   len(users),
		"ActiveUsers": activeUsers,
		"ServerCount": len(servers),
	})
}

func (h *Handler) AdminListUsers(c *fiber.Ctx) error {
	type userView struct {
		ID        string
		Name      string
		UUID      string
		Active    bool
		Expired   bool
		ExpiresAt *time.Time
		ServerIDs []string
	}

	users := h.cfg.AllUsers()
	views := make([]userView, len(users))
	for i, u := range users {
		views[i] = userView{
			ID:        u.ID,
			Name:      u.Name,
			UUID:      u.UUID,
			Active:    u.Active,
			Expired:   u.IsExpired(),
			ExpiresAt: u.ExpiresAt,
			ServerIDs: u.ServerIDs,
		}
	}
	return c.Render("users", fiber.Map{"Users": views})
}

func (h *Handler) AdminNewUserForm(c *fiber.Ctx) error {
	servers := h.cfg.AllServers()
	return c.Render("user_form", fiber.Map{"Servers": servers})
}

func (h *Handler) AdminCreateUser(c *fiber.Ctx) error {
	u := model.User{
		ID:     fmt.Sprintf("usr_%s", model.GenerateToken()[4:16]),
		Name:   c.FormValue("name"),
		UUID:   model.GenerateUUID(),
		Token:  model.GenerateToken(),
		Active: true,
	}

	if raw := c.Context().PostArgs().PeekMulti("server_ids"); len(raw) > 0 {
		u.ServerIDs = make([]string, len(raw))
		for i, v := range raw {
			u.ServerIDs[i] = string(v)
		}
	}

	if err := u.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if exp := c.FormValue("expires_at"); exp != "" {
		t, err := parseTime(exp)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid expires_at format")
		}
		u.ExpiresAt = &t
	}

	h.cfg.AddUser(u)
	if err := h.cfg.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("save failed")
	}

	return c.Redirect("/admin/users")
}

func (h *Handler) AdminDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	h.cfg.RemoveUser(id)
	h.cfg.Save()
	return c.Redirect("/admin/users")
}

func (h *Handler) AdminListServers(c *fiber.Ctx) error {
	return c.Render("servers", fiber.Map{"Servers": h.cfg.AllServers()})
}

func (h *Handler) AdminNewServerForm(c *fiber.Ctx) error {
	return c.Render("server_form", fiber.Map{})
}

func (h *Handler) AdminCreateServer(c *fiber.Ctx) error {
	s := model.Server{
		ID:          fmt.Sprintf("srv_%d", time.Now().UnixMilli()),
		Name:        c.FormValue("name"),
		Host:        c.FormValue("host"),
		Security:    c.FormValue("security"),
		Transport:   c.FormValue("transport"),
		SNI:         c.FormValue("sni"),
		Fingerprint: c.FormValue("fingerprint"),
		Flow:        c.FormValue("flow"),
		Path:        c.FormValue("path"),
		HostHeader:  c.FormValue("host_header"),
		PBK:         c.FormValue("pbk"),
		SID:         c.FormValue("sid"),
		ServiceName: c.FormValue("service_name"),
	}

	port := c.FormValue("port")
	if port != "" {
		fmt.Sscanf(port, "%d", &s.Port)
	}

	if err := s.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	h.cfg.AddServer(s)
	if err := h.cfg.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("save failed")
	}

	return c.Redirect("/admin/servers")
}

func (h *Handler) AdminDeleteServer(c *fiber.Ctx) error {
	id := c.Params("id")
	h.cfg.RemoveServer(id)
	h.cfg.Save()
	return c.Redirect("/admin/servers")
}
