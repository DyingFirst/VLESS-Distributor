package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"

	"github.com/dyingfirst/subscribe-server/internal/model"
)

type AppConfig struct {
	Listen     string `yaml:"listen"`
	APIKey     string `yaml:"api_key"`
	ConfigPath string `yaml:"config_path"`
}

type Config struct {
	App     AppConfig      `yaml:"app"`
	Servers []model.Server `yaml:"servers"`
	Users   []model.User   `yaml:"users"`

	mu sync.RWMutex
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.App.ConfigPath == "" {
		cfg.App.ConfigPath = path
	}
	if cfg.App.Listen == "" {
		cfg.App.Listen = ":8080"
	}

	return cfg, nil
}

func (c *Config) Save() error {
	c.mu.RLock()
	data, err := yaml.Marshal(c)
	c.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(c.App.ConfigPath, data, 0644)
}

func (c *Config) FindUserByToken(token string) *model.User {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := range c.Users {
		if c.Users[i].Token == token {
			return &c.Users[i]
		}
	}
	return nil
}

func (c *Config) FindUserByID(id string) *model.User {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := range c.Users {
		if c.Users[i].ID == id {
			return &c.Users[i]
		}
	}
	return nil
}

func (c *Config) FindServerByID(id string) *model.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := range c.Servers {
		if c.Servers[i].ID == id {
			return &c.Servers[i]
		}
	}
	return nil
}

func (c *Config) GetServersForUser(u *model.User) []model.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(u.ServerIDs) == 0 {
		result := make([]model.Server, len(c.Servers))
		copy(result, c.Servers)
		return result
	}

	var result []model.Server
	for _, srv := range c.Servers {
		for _, sid := range u.ServerIDs {
			if srv.ID == sid {
				result = append(result, srv)
				break
			}
		}
	}
	return result
}

func (c *Config) AddUser(u model.User) {
	c.mu.Lock()
	c.Users = append(c.Users, u)
	c.mu.Unlock()
}

func (c *Config) RemoveUser(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, u := range c.Users {
		if u.ID == id {
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Config) AddServer(s model.Server) {
	c.mu.Lock()
	c.Servers = append(c.Servers, s)
	c.mu.Unlock()
}

func (c *Config) RemoveServer(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, s := range c.Servers {
		if s.ID == id {
			c.Servers = append(c.Servers[:i], c.Servers[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Config) AllUsers() []model.User {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]model.User, len(c.Users))
	copy(result, c.Users)
	return result
}

func (c *Config) AllServers() []model.Server {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]model.Server, len(c.Servers))
	copy(result, c.Servers)
	return result
}

func (c *Config) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("fsnotify: %v, hot-reload disabled", err)
		return
	}

	if err := watcher.Add(c.App.ConfigPath); err != nil {
		log.Printf("fsnotify add: %v, hot-reload disabled", err)
		watcher.Close()
		return
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					newCfg, err := Load(c.App.ConfigPath)
					if err != nil {
						log.Printf("config reload failed: %v", err)
						continue
					}
					c.mu.Lock()
					c.App = newCfg.App
					c.Servers = newCfg.Servers
					c.Users = newCfg.Users
					c.mu.Unlock()
					log.Println("config reloaded")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("fsnotify error: %v", err)
			}
		}
	}()

	log.Println("watching config for changes")
}
