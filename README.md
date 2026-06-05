# Subscribe Server

A lightweight self-hosted service that serves VLESS client configurations via subscription links, with a built-in admin panel.

## Features

- Subscription endpoint returning Base64-encoded VLESS URI lists
- Support for TLS, Reality, WebSocket, gRPC, HTTPUpgrade transports
- Admin panel with dashboard, user and server management
- JSON API for programmatic access
- YAML-based configuration with hot-reload
- Single binary with embedded templates, no external dependencies

## Quick Start

```bash
# Build
go build ./cmd/subscribe-server

# Create config from example
cp config.example.yaml config.yaml

# Run
./subscribe-server -config config.yaml
```

The server listens on `:8080` by default.

## Docker

Multi-platform image (`linux/amd64`, `linux/arm64`) is available on Docker Hub:

```bash
cp config.example.yaml config.yaml
# Edit config.yaml with your real values

docker compose up -d
```

Or run directly:

```bash
docker run -d -p 8080:8080 -v ./config.yaml:/data/config.yaml:rw dyingfirst/vless-distributor:latest
```

## Configuration

Configuration is stored in a single YAML file. See `config.example.yaml` for a full reference.

### Application Settings

```yaml
app:
  listen: ":8080"         # Listen address
  api_key: "changeme123"  # API key for management endpoints and admin login
  config_path: "config.yaml"
```

### Server Definition

```yaml
servers:
  - id: "srv_1"
    name: "US East"           # Display name
    host: "us1.example.com"
    port: 443
    security: "tls"           # none | tls | reality
    transport: "ws"           # tcp | ws | grpc | httpupgrade
    path: "/ws"               # For ws/httpupgrade
    host_header: "us1.example.com"
    sni: "us1.example.com"
    fingerprint: "chrome"
    flow: ""                  # xtls-rprx-vision (for reality+tcp)
    pbk: ""                   # Reality public key
    sid: ""                   # Reality short ID
    service_name: ""          # gRPC service name
```

### User Definition

```yaml
users:
  - id: "usr_1"
    name: "alice"
    uuid: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
    token: "sub_example_token_alice"  # Subscription token (URL-safe)
    active: true
    expires_at: "2026-12-31T23:59:59Z"  # RFC3339, omit for never
    server_ids: []  # Empty = all servers; otherwise specific server IDs
```

UUID and token are auto-generated when creating users via the API or admin panel.

## Endpoints

### Subscription (public)

```
GET /sub/{token}
```

Returns a Base64-encoded list of VLESS URIs for the user identified by the token. Use this URL as the subscription link in any VLESS-compatible client (v2rayN, v2rayNG, Streisand, etc.).

Example:

```bash
curl http://localhost:8080/sub/sub_example_token_alice
```

### Admin Panel

Access `/admin/` in a browser. Log in using the `api_key` value from config.

| Route | Description |
|---|---|
| `GET /admin/` | Dashboard with user/server counts |
| `GET /admin/users` | User list |
| `GET /admin/users/new` | Create user form |
| `POST /admin/users` | Create user |
| `POST /admin/users/{id}/delete` | Delete user |
| `GET /admin/servers` | Server list |
| `GET /admin/servers/new` | Create server form |
| `POST /admin/servers` | Create server |
| `POST /admin/servers/{id}/delete` | Delete server |
| `GET /admin/login` | Login page |
| `POST /admin/login` | Authenticate |
| `POST /admin/logout` | Logout |

### JSON API

All API endpoints require authentication via `Authorization: Bearer <api_key>` header or `?api_key=` query parameter.

| Method | Route | Description |
|---|---|---|
| `GET` | `/api/v1/users` | List all users |
| `POST` | `/api/v1/users` | Create user |
| `DELETE` | `/api/v1/users/{id}` | Delete user |
| `GET` | `/api/v1/servers` | List all servers |
| `POST` | `/api/v1/servers` | Create server |
| `DELETE` | `/api/v1/servers/{id}` | Delete server |

Create user (JSON body):

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer changeme123" \
  -H "Content-Type: application/json" \
  -d '{"name": "alice", "expires_at": "2026-12-31T23:59:59Z"}'
```

Create server (JSON body):

```bash
curl -X POST http://localhost:8080/api/v1/servers \
  -H "Authorization: Bearer changeme123" \
  -H "Content-Type: application/json" \
  -d '{"name": "US East", "host": "us1.example.com", "port": 443, "security": "tls", "transport": "ws", "path": "/ws", "sni": "us1.example.com"}'
```

## Hot Reload

The server watches the config file for changes. Editing `config.yaml` while the server is running will automatically reload the configuration without a restart.
