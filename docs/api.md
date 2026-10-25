# HTTP API reference

The API listens on port `8081` by default (`--api-port`). The dashboard listens
on port `8080` by default (`--port`). Both servers expose `GET /health` and
return `{"status":"ok"}`; use the port to distinguish them.

Box and event routes use the `/api/v1` prefix and JSON request and response
bodies, except `DELETE` on a box, which returns no body on success.

## Endpoints

| Method | Path | Success response | Description |
| :---: | :--- | :--- | :--- |
| `GET` | `/health` | `200`, `{"status":"ok"}` | API health check |
| `GET` | `/api/v1/boxes` | `200`, array of boxes | List boxes |
| `POST` | `/api/v1/boxes` | `201`, created box | Create a box |
| `GET` | `/api/v1/boxes/{id}` | `200`, box | Get a box |
| `PUT` | `/api/v1/boxes/{id}` | `200`, replaced box; `201`, new box | Replace or create a box |
| `DELETE` | `/api/v1/boxes/{id}` | `204`, no body | Delete a box |
| `POST` | `/api/v1/boxes/{id}/events` | `201`, accepted event | Update a box's status |

Example health response:

```bash
curl -i http://localhost:8081/health
```

```json
{"status":"ok"}
```

The older singular `/api/v1/box...` routes remain available and return a
`Warning` response header. Use the plural routes above for new integrations.

## Box schema

Fields accepted when creating or replacing a box:

| Field | Type | Required | Meaning |
| :--- | :--- | :---: | :--- |
| `id` | string | No for create; yes for `PUT` | Unique identifier. Create generates one if omitted. |
| `name` | string | No | Box name. |
| `displayName` | string | No | Optional tile label. |
| `description` | string | No | Description shown in the box details. |
| `parent` | string | No | Parent box ID. |
| `size` | string | No | One of `dot`, `micro`, `dmicro`, `small`, `dsmall`, `medium`, `dmedium`, `large`, `dlarge`, `xlarge`. |
| `status` | string | No | One of `grey`, `gray`, `red`, `amber`, `green`, `noUpdate`. `gray` is accepted as an alias and is returned as `grey`. |
| `maxTBU` | duration | No | Mark `noUpdate` when no update arrives within this duration. |
| `expireAfter` | duration | No | Delete the box after this duration without an update. |
| `lastMessage` | string | No | Initial status message. |
| `links` | array of `{ "name": string, "url": string }` | No | Links shown in box details. |
| `info` | object of string values | No | Additional details shown in the dashboard. |
| `messages` | array of message objects | No | Existing event history; normally omit when creating a box. |

Responses also include `lastUpdate`. Each `messages` entry has `message`,
`status`, and `timeStamp`, and may include `maxTBU` and `expireAfter` when those
values were supplied with that event. The server sets `lastUpdate`; it also
creates an initial history entry when a new box is created without `messages`.
Omitted `size` and `status` values use their zero values, `dot` and `grey`.

Example box response (also returned by create, get, and replace):

```json
{
  "id": "my-service",
  "name": "My Service",
  "size": "medium",
  "status": "green",
  "messages": [{"message":"All checks passed","status":"green","timeStamp":"2026-10-25T12:00:00Z"}],
  "lastUpdate": "2026-10-25T12:00:00Z",
  "lastMessage": "All checks passed",
  "links": []
}
```

Durations should be strings accepted by Go's duration format, such as `"30m"`,
`"6h"`, or `"1h30m"`. A plain numeric string is also interpreted as seconds.
Use `"0s"` in an event to disable `maxTBU` or `expireAfter`.

## Create a box

```bash
curl -i -X POST http://localhost:8081/api/v1/boxes \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "my-service",
    "name": "My Service",
    "size": "medium",
    "status": "grey",
    "maxTBU": "6h",
    "links": [{"name": "Runbook", "url": "https://example.com/runbook"}],
    "info": {"team": "platform"}
  }'
```

Returns `201 Created` and the created box as JSON. If `id` is omitted, the
server generates one. An existing ID currently produces `500 Internal Server
Error`.

## List and get boxes

```bash
curl -i http://localhost:8081/api/v1/boxes
curl -i http://localhost:8081/api/v1/boxes/my-service
```

Both return `200 OK`. The list response is a JSON array; an unknown ID returns
`404 Not Found` with an error response.

The list response contains box objects in the shape shown above; the get
response contains one box object.

## Replace or create a box

```bash
curl -i -X PUT http://localhost:8081/api/v1/boxes/my-service \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "my-service",
    "name": "My Service",
    "size": "large",
    "status": "green"
  }'
```

The handler requires `id` in the JSON body and uses that body value to choose
which box to replace or create. It currently does not compare the body ID with
the `{id}` in the URL, so keep them identical. Returns `200 OK` when replacing
an existing box and `201 Created` when creating one. A missing body ID returns
`400 Bad Request`.

## Post a status event

```bash
curl -i -X POST http://localhost:8081/api/v1/boxes/my-service/events \
  -H 'Content-Type: application/json' \
  -d '{
    "status": "green",
    "lastMessage": "All checks passed",
    "maxTBU": "30m"
  }'
```

The event body accepts `status`, `lastMessage`, `maxTBU`, and `expireAfter`.
The URL supplies the box ID; any `id` in the body is ignored. Omitted timing
fields leave the current box settings unchanged. Supplying `"0s"` clears that
setting. The `201 Created` response contains the event with its `type` set to
`updateBox`.

Example event response:

```json
{"status":"green","lastMessage":"All checks passed","maxTBU":"30m","expireAfter":null,"type":"updateBox"}
```

An unknown box returns `404 Not Found`; malformed JSON returns `400 Bad
Request`.

## Delete a box

```bash
curl -i -X DELETE http://localhost:8081/api/v1/boxes/my-service
```

Returns `204 No Content` when deleted, or `404 Not Found` with
`{"error":"box not found"}` if no box has that ID.

## Errors

Most API errors use this JSON shape:

```json
{"message":"could not decode data received","error":"details"}
```

`error` is omitted on some server errors. The delete-not-found response is a
separate shape containing only `error`. Current status codes by operation:

| Operation | Error status codes |
| :--- | :--- |
| Create box | `400` invalid JSON; `500` creation failure, including duplicate ID |
| Get box | `404` unknown ID |
| Replace/create box | `400` invalid JSON or missing body ID; `500` replacement/creation failure |
| Post event | `400` invalid JSON; `404` unknown box; `500` update failure |
| Delete box | `404` unknown ID |

The API does not currently publish a formal schema or interactive API page.
