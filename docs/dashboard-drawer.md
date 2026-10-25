# Dashboard drawer behaviour

Clicking a tile opens a drawer with Overview and Events tabs. Overview contains
the current box settings and its latest update; Events contains its history.

## Layout fixes

The dashboard grows with the tile grid instead of reserving nearly a full
viewport below the status bar, which previously caused unnecessary scrolling.
Scrollbar styling applies to the document as well as the drawer, while wheel
and touch scrolling remain available when content overflows.

The drawer content fills the space below its header and tabs. Its content area
handles scrolling, and the event log has no fixed height limit or nested scroll
area. Event cards retain the coloured rail; the repeated colour name is visually
hidden and remains available to screen readers. Empty event headings collapse
when no timing indicators are present.

## Why box creation changes in Go

Both the server-rendered drawer and the JavaScript that handles live box
creation read event indicators from `Box.Messages`. Previously, creating a box
with `maxTBU` or `expireAfter` set could leave `Messages` empty. Status Details
could display the settings, but Latest update and Event log had no creation
entry from which to display those values.

`addBox` in `internal/server/boxes.go` now creates one initial history entry
when no history is supplied. It uses the box's initial status, message, timing
settings, and creation timestamp. Creating the entry before storing and
broadcasting the box gives live clients and reloaded pages the same data.
No additional event API request is required.

This is an observable backend change supporting the drawer: GET responses,
the `createBox` SSE payload, and persisted data include the initial entry. It
occupies one of the existing 30 history slots until newer events displace it.
Supplied history is preserved, and existing persisted history is not backfilled.
The status-update logic in Go retains its original behaviour.

## Current settings and event fields

Status Details shows the box's current `maxTBU` and `expireAfter`. Omitting a
field from a status update leaves that setting unchanged. Explicitly sending
`"0s"` disables the setting and removes it from Status Details.

Latest update and each Event log entry show only timing fields supplied with
that request. They do not inherit timing values from the box's current state.
The initial creation entry includes any enabled timing settings supplied when
the box was created.

For a box created with `maxTBU: "30m"` and `expireAfter: "1h"`:

| Subsequent event payload | Status Details afterwards | Timing indicators on that event |
| --- | --- | --- |
| Status and message only | Max TBU 30m; expiry 1h | None |
| Adds `maxTBU: "45m"` | Max TBU 45m; expiry 1h | Max TBU 45m only |
| Adds `expireAfter: "0s"` | Max TBU 45m; expiry disabled | Expiry 0s only |

The browser ignores omitted or null timing fields when updating Status Details.
The API sends omitted event fields as null in SSE messages, so treating null as
a removal would incorrectly erase the drawer's current settings. Event cards
still omit those null fields. Explicit zero durations remain visible in event
history as the request that disabled the setting.

## Verification

`TestBoxEventMetadata` in the existing `boxstore_test.go` covers the creation
entry, partial updates, explicit clears, supplied-history preservation, and
the distinction between stored current settings and event-specific metadata.
Run the Go suite with `go test ./...`.

Browser checks cover initial creation, status-only updates, partial timing
changes, explicit clears, reopening the drawer, and page reloads. Layout checks
cover unnecessary page scrolling, hidden scrollbars, drawer height, and access
to overflowing content. The stylesheet also validates without diagnostics using
the installed VS Code CSS language server. Safari and Firefox were not checked.
