# alive

A real-time status dashboard. Track the health of services, cron jobs,
or anything you can script a health check for — displayed as a live, auto-updating
grid of coloured tiles in a browser.

## How it works

Everything on the dashboard is a **box** — a coloured tile representing one thing
you want to monitor. External scripts post status updates to the REST API;
the dashboard updates instantly in any open browser via server-sent events (SSE),
with no polling or page refresh required.

### Box statuses

| Status | Meaning |
| :----: | :------ |
| `grey` | Unknown |
| `green` | OK |
| `amber` | Warning |
| `red` | Error |
| `noUpdate` | No update received within the configured time window |

### Box sizes

Boxes can be sized from smallest to largest:
`dot`, `micro`, `dmicro`, `small`, `dsmall`, `medium`, `dmedium`, `large`, `dlarge`, `xlarge`.

### Box details

Click a tile to open its details drawer. Status Details shows the box's current
settings; Latest update and Event log show only timing fields supplied with
each request. See [drawer behaviour](docs/dashboard-drawer.md) for examples,
the creation-history rationale, and verification notes.

## Running

```bash
alive [OPTIONS]
```

| Flag | Default | Description |
| :--: | :-----: | :---------- |
| `--port` / `-p` | `8080` | Dashboard port |
| `--api-port` | `8081` | API port |
| `--data-path` / `-d` | `$HOME/.alive/data` | Where box state is persisted |
| `--static-path` | `$HOME/.alive/static` | Where static files are served from |
| `--default-static` | | Use built-in CSS/JS instead of files on disk |
| `--run-demo` | | Run a self-contained demo using a temporary directory |
| `--debug` | | Enable debug logging |

### Docker

Build the image without reusing Docker build layers:

```bash
docker build --no-cache -t alive:latest .
```

Run it with a named volume for persisted box data. The `--rm` flag removes the
container automatically when it stops; it does not remove the image or the
named data volume.

```bash
docker run --rm --name alive \
  -p 8080:8080 \
  -p 8081:8081 \
  -v alive-data:/data \
  alive:latest
```

The Docker image enables `--default-static` by default, so the CSS and
JavaScript embedded during the image build are copied into the container's
static directory and current frontend changes are used after rebuilding.

Check the running container and endpoints:

```bash
docker ps
curl http://localhost:8081/health
curl http://localhost:8081/api/v1/boxes
```

To stop and remove an existing container before recreating it:

```bash
docker rm -f alive
```

To remove stopped containers and unused images later:

```bash
docker container prune
docker image prune
```

### Demo mode

```bash
alive --run-demo
```

Starts with example boxes pre-populated, using a temporary directory — no setup required.

## API

See the [API reference](docs/api.md) for endpoint behavior, schemas, response
codes, and request examples. The dashboard and API each expose `/health` on
their own ports: dashboard `8080`, API `8081` by default.

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

Box state is saved to disk every minute and on shutdown. On startup,
state is restored from the data file so boxes survive restarts.

## Examples

The `examples/` directory contains shell scripts showing common usage patterns
including SSL certificate checks, DNS tests, connectivity checks, and ad-hoc job monitoring.
