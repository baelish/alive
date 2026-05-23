# alive

A real-time status dashboard. Track the health of services, cron jobs, or anything you can script a health check for — displayed as a live, auto-updating grid of coloured tiles in a browser.

## How it works

Everything on the dashboard is a **box** — a coloured tile representing one thing you want to monitor. External scripts post status updates to the REST API; the dashboard updates instantly in any open browser via server-sent events (SSE), with no polling or page refresh required.

### Box statuses

| Status | Meaning |
|--------|---------|
| `grey` | Unknown |
| `green` | OK |
| `amber` | Warning |
| `red` | Error |
| `noUpdate` | No update received within the configured time window |

### Box sizes

Boxes can be sized from smallest to largest: `dot`, `micro`, `dmicro`, `small`, `dsmall`, `medium`, `dmedium`, `large`, `dlarge`, `xlarge`.

## Running

```
alive [OPTIONS]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--port` / `-p` | `8080` | Dashboard port |
| `--api-port` | `8081` | API port |
| `--data-path` / `-d` | `$HOME/.alive/data` | Where box state is persisted |
| `--static-path` | `$HOME/.alive/static` | Where static files are served from |
| `--default-static` | | Use built-in CSS/JS instead of files on disk |
| `--run-demo` | | Run a self-contained demo using a temporary directory |
| `--debug` | | Enable debug logging |

### Docker

```
docker build -t alive .
docker run -p 8080:8080 -p 8081:8081 alive
```

### Demo mode

```
alive --run-demo
```

Starts with example boxes pre-populated, using a temporary directory — no setup required.

## API

The API listens on port `8081` by default.

### Boxes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/boxes` | List all boxes |
| `POST` | `/api/v1/boxes` | Create a box |
| `GET` | `/api/v1/boxes/{id}` | Get a specific box |
| `PUT` | `/api/v1/boxes/{id}` | Replace a box (creates if not found) |
| `DELETE` | `/api/v1/boxes/{id}` | Delete a box |
| `POST` | `/api/v1/boxes/{id}/events` | Post a status update to a box |
| `GET` | `/health` | Health check |

### Create a box

```bash
curl -X POST http://localhost:8081/api/v1/boxes \
  -H "Content-Type: application/json" \
  -d '{
    "id": "my-service",
    "name": "My Service",
    "size": "medium",
    "status": "grey",
    "maxTBU": "6h"
  }'
```

Box fields:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique ID (auto-generated if omitted) |
| `name` | string | Display name |
| `displayName` | string | Alternative display name shown on the tile |
| `description` | string | Shown on the detail page |
| `size` | string | Tile size (see sizes above) |
| `status` | string | Initial status |
| `maxTBU` | duration | Flip to `noUpdate` if no event arrives within this window (e.g. `"6h"`, `"30m"`) |
| `expireAfter` | duration | Auto-delete the box after this duration without an update |
| `links` | array | `[{"name": "...", "url": "..."}]` — shown on the detail page |
| `info` | object | Arbitrary key/value pairs shown on the detail page |

### Post a status update

```bash
curl -X POST http://localhost:8081/api/v1/boxes/my-service/events \
  -H "Content-Type: application/json" \
  -d '{
    "status": "green",
    "lastMessage": "All checks passed"
  }'
```

## Go client

A Go client package is included:

```go
import "github.com/baelish/alive/client"

c := client.NewClient("http://localhost:8081")
c.CreateBox(box)
c.GetAllBoxes()
c.GetBox("my-service")
c.ReplaceBox(box)
c.DeleteBox("my-service")
```

## State persistence

Box state is saved to disk every minute and on shutdown. On startup, state is restored from the data file so boxes survive restarts.

## Examples

The `examples/` directory contains shell scripts showing common usage patterns including SSL certificate checks, DNS tests, connectivity checks, and ad-hoc job monitoring.
