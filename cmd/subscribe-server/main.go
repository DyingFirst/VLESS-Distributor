package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/dyingfirst/subscribe-server/internal/config"
	"github.com/dyingfirst/subscribe-server/internal/handler"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg.Watch()

	sub, err := fs.Sub(templateFS, "templates")
	if err != nil {
		log.Fatalf("templates: %v", err)
	}
	engine := html.NewFileSystem(http.FS(sub), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	h := handler.New(cfg)

	// Public
	app.Get("/sub/:token", h.GetSubscription)

	// JSON API
	api := app.Group("/api/v1", h.APIKeyMiddleware)
	api.Get("/users", h.ListUsers)
	api.Post("/users", h.CreateUser)
	api.Delete("/users/:id", h.DeleteUser)
	api.Get("/servers", h.ListServers)
	api.Post("/servers", h.CreateServer)
	api.Delete("/servers/:id", h.DeleteServer)

	// Admin login (no auth)
	app.Get("/admin/login", h.AdminLoginPage)
	app.Post("/admin/login", h.AdminLogin)
	app.Post("/admin/logout", h.AdminLogout)

	// Admin panel (cookie auth)
	admin := app.Group("/admin", h.CookieAuthMiddleware)
	admin.Get("/", h.AdminDashboard)
	admin.Get("/users", h.AdminListUsers)
	admin.Get("/users/new", h.AdminNewUserForm)
	admin.Post("/users", h.AdminCreateUser)
	admin.Post("/users/:id/delete", h.AdminDeleteUser)
	admin.Get("/servers", h.AdminListServers)
	admin.Get("/servers/new", h.AdminNewServerForm)
	admin.Post("/servers", h.AdminCreateServer)
	admin.Post("/servers/:id/delete", h.AdminDeleteServer)

	addr := cfg.App.Listen
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
