package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Broker which will be created in this program. It is responsible
// for keeping a list of which clients (browsers) are currently attached
// and broadcasting events (messages) to those clients.
type Broker struct {

	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients.  (The values
	// are just booleans and are meaningless.)
	clients map[chan string]bool

	// Channel into which new clients can be pushed
	newClients chan chan string

	// Channel into which disconnected clients should be pushed
	defunctClients chan chan string

	// Channel into which messages are pushed to be broadcast out
	// to attached clients.
	messages chan string

	// Channel to query client count (for testing)
	clientCount chan int
}

// Start method, this Broker method starts a new goroutine.  It handles
// the addition & removal of clients, as well as the broadcasting
// of messages out to clients that are currently attached.
func (b *Broker) Start(ctx context.Context) {

	// Start a goroutine
	go func() {

		// Loop endlessly
		for {

			// Block until we receive from one of the
			// three following channels.
			select {

			case <-ctx.Done():
				return

			case s := <-b.newClients:

				// There is a new client attached and we
				// want to start sending them messages.
				b.clients[s] = true
				logger.Info("Added new client", zap.Int("currentClientCount", len(b.clients)))

			case s := <-b.defunctClients:

				// A client has detached and we want to
				// stop sending them messages.
				delete(b.clients, s)
				close(s)

				logger.Info("Removed client", zap.Int("currentClientCount", len(b.clients)))

			case b.clientCount <- len(b.clients):
				// Respond to client count query (non-blocking send from caller's perspective)

			case msg := <-b.messages:

				// There is a new message to send.  For each
				// attached client, push the new message
				// into the client's message channel.
				for s := range b.clients {
					// Non-blocking send to prevent slow clients from blocking broker
					select {
					case s <- msg:
						// Message sent successfully
					default:
						// Client's buffer is full, drop the message
						// This prevents one slow client from blocking all others
						logger.Warn("Dropped message for slow client")
					}
				}
			}
		}
	}()
}

// ClientCount returns the current number of connected clients (thread-safe)
func (b *Broker) ClientCount() int {
	return <-b.clientCount
}

// This Broker method handles and HTTP request at the "/events/" URL.
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel, over which the broker can
	// send this client messages.
	// Buffered to prevent slow clients from blocking the broker
	messageChan := make(chan string, 100)

	// Add this client to the map of those that should
	// receive updates
	b.newClients <- messageChan

	// Listen to the closing of the http connection via the request context
	// The context is cancelled when the client disconnects
	go func() {
		<-r.Context().Done()
		// Remove this client from the map of attached clients
		// when the client disconnects
		b.defunctClients <- messageChan
		logger.Warn("http connection just closed")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Don't close the connection, instead loop endlessly.
	for {

		// Read from our messageChan.
		msg, open := <-messageChan

		if !open {
			// If our messageChan was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		fmt.Fprintf(w, "data: %s\n\n", msg)

		// Flush the response.  This is only possible if
		// the repsonse supports streaming.
		f.Flush()
	}

	// Done.
	logger.Info("Finished HTTP request", zap.String("path", r.URL.Path))
}

// Send keepalives to the status bar.
func runKeepalives(ctx context.Context) {
	if options.Debug {
		logger.Info("Starting keepalive routine")
	}
	// Generate a regular keepalive message that gets pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.
	for {
		select {
		case <-ctx.Done():
			if options.Debug {
				logger.Info("Stopping keepalive routine")
			}
			return

		case <-time.After(time.Duration(3 * time.Second)):
		}

		// Send a keepalive
		events.messages <- `{"type": "keepalive"}`
	}
}

// Main routine
func runSSE(ctx context.Context) (b *Broker) {
	if options.Debug {
		logger.Info("Starting SSE broker")
	}

	// Make a new Broker instance
	b = &Broker{
		clients:        make(map[chan string]bool),
		newClients:     make(chan (chan string)),
		defunctClients: make(chan (chan string)),
		messages:       make(chan string),
		clientCount:    make(chan int),
	}

	// Start processing events
	b.Start(ctx)

	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	http.Handle("/events/", b)

	return b
}
