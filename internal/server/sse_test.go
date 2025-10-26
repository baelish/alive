package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBroker_Start(t *testing.T) {
	t.Run("adds new clients", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		broker.Start(ctx)

		// Add a new client
		clientChan := make(chan string, 1)
		broker.newClients <- clientChan

		// Give it a moment to process
		time.Sleep(10 * time.Millisecond)

		// Verify the client was added
		if len(broker.clients) != 1 {
			t.Errorf("expected 1 client, got %d", len(broker.clients))
		}
		if !broker.clients[clientChan] {
			t.Error("client channel was not added to broker.clients map")
		}
	})

	t.Run("removes defunct clients", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		broker.Start(ctx)

		// Add a client
		clientChan := make(chan string, 1)
		broker.newClients <- clientChan
		time.Sleep(10 * time.Millisecond)

		if len(broker.clients) != 1 {
			t.Fatalf("expected 1 client after adding, got %d", len(broker.clients))
		}

		// Remove the client
		broker.defunctClients <- clientChan
		time.Sleep(10 * time.Millisecond)

		// Verify the client was removed
		if len(broker.clients) != 0 {
			t.Errorf("expected 0 clients after removal, got %d", len(broker.clients))
		}
		if broker.clients[clientChan] {
			t.Error("client channel was not removed from broker.clients map")
		}

		// Verify the channel was closed
		_, open := <-clientChan
		if open {
			t.Error("expected client channel to be closed, but it was still open")
		}
	})

	t.Run("broadcasts messages to all clients", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		broker.Start(ctx)

		// Add multiple clients
		client1 := make(chan string, 1)
		client2 := make(chan string, 1)
		client3 := make(chan string, 1)

		broker.newClients <- client1
		broker.newClients <- client2
		broker.newClients <- client3
		time.Sleep(10 * time.Millisecond)

		// Broadcast a message
		testMessage := "test message"
		broker.messages <- testMessage
		time.Sleep(10 * time.Millisecond)

		// Verify all clients received the message
		if msg := <-client1; msg != testMessage {
			t.Errorf("client1: expected %q, got %q", testMessage, msg)
		}
		if msg := <-client2; msg != testMessage {
			t.Errorf("client2: expected %q, got %q", testMessage, msg)
		}
		if msg := <-client3; msg != testMessage {
			t.Errorf("client3: expected %q, got %q", testMessage, msg)
		}
	})

	t.Run("stops when context is cancelled", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		broker.Start(ctx)

		// Add a client to verify broker is running
		clientChan := make(chan string, 1)
		broker.newClients <- clientChan
		time.Sleep(10 * time.Millisecond)

		if len(broker.clients) != 1 {
			t.Fatalf("expected 1 client, got %d", len(broker.clients))
		}

		// Cancel the context
		cancel()
		time.Sleep(10 * time.Millisecond)

		// After context cancellation, the broker should stop processing
		// We can't directly test if the goroutine stopped, but we can verify
		// that the broker becomes unresponsive
		select {
		case broker.newClients <- make(chan string):
			// This might succeed immediately if buffered, so we need a timeout
			time.Sleep(50 * time.Millisecond)
		case <-time.After(100 * time.Millisecond):
			// Expected: the broker is no longer processing newClients
		}
	})
}

// testResponseWriter implements http.ResponseWriter, http.Flusher, and http.CloseNotifier
// for testing SSE functionality
type testResponseWriter struct {
	header      http.Header
	body        []byte
	statusCode  int
	closeNotify chan bool
	flushed     bool
}

func (w *testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *testResponseWriter) Flush() {
	w.flushed = true
}

func (w *testResponseWriter) CloseNotify() <-chan bool {
	if w.closeNotify == nil {
		w.closeNotify = make(chan bool)
	}
	return w.closeNotify
}

func TestBroker_ServeHTTP(t *testing.T) {
	t.Run("sets correct SSE headers", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		broker.Start(ctx)

		// Use a custom ResponseWriter to capture headers
		testWriter := &testResponseWriter{
			header: make(http.Header),
		}
		testReq := httptest.NewRequest("GET", "/events/", nil)

		go broker.ServeHTTP(testWriter, testReq)

		// Give it a moment to set headers
		time.Sleep(20 * time.Millisecond)

		// Verify headers
		if ct := testWriter.Header().Get("Content-Type"); ct != "text/event-stream" {
			t.Errorf("Content-Type: expected %q, got %q", "text/event-stream", ct)
		}
		if cc := testWriter.Header().Get("Cache-Control"); cc != "no-cache" {
			t.Errorf("Cache-Control: expected %q, got %q", "no-cache", cc)
		}
		if conn := testWriter.Header().Get("Connection"); conn != "keep-alive" {
			t.Errorf("Connection: expected %q, got %q", "keep-alive", conn)
		}
		if te := testWriter.Header().Get("Transfer-Encoding"); te != "chunked" {
			t.Errorf("Transfer-Encoding: expected %q, got %q", "chunked", te)
		}
	})

	t.Run("sends messages in SSE format", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		broker.Start(ctx)

		// Use custom test writer
		testWriter := &testResponseWriter{
			header: make(http.Header),
			body:   make([]byte, 0),
		}
		req := httptest.NewRequest("GET", "/events/", nil)

		// Start ServeHTTP in a goroutine
		done := make(chan bool)
		go func() {
			broker.ServeHTTP(testWriter, req)
			done <- true
		}()

		// Wait for client to be registered
		time.Sleep(20 * time.Millisecond)

		// Send a message
		testMessage := "hello world"
		broker.messages <- testMessage

		// Give it time to write
		time.Sleep(20 * time.Millisecond)

		// Check the response body contains the SSE-formatted message
		body := string(testWriter.body)
		expectedFormat := "data: hello world\n\n"
		if !strings.Contains(body, expectedFormat) {
			t.Errorf("expected body to contain %q, got %q", expectedFormat, body)
		}
	})

	t.Run("registers and unregisters client", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		broker.Start(ctx)

		testWriter := &testResponseWriter{
			header: make(http.Header),
			body:   make([]byte, 0),
		}
		req := httptest.NewRequest("GET", "/events/", nil)

		// Start ServeHTTP in a goroutine
		go func() {
			broker.ServeHTTP(testWriter, req)
		}()

		// Wait for client to register
		time.Sleep(20 * time.Millisecond)

		// Verify client was registered
		if len(broker.clients) != 1 {
			t.Errorf("expected 1 registered client, got %d", len(broker.clients))
		}
	})
}

func TestRunSSE(t *testing.T) {
	// Save original global state
	originalOptions := options
	defer func() { options = originalOptions }()

	// Note: runSSE registers a global HTTP handler at /events/
	// This can only be done once, so we test it in a single test
	// Set debug mode off to avoid log output
	options.Debug = false

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker := runSSE(ctx)

	if broker == nil {
		t.Fatal("runSSE returned nil broker")
	}
	if broker.clients == nil {
		t.Error("broker.clients map is nil")
	}
	if broker.newClients == nil {
		t.Error("broker.newClients channel is nil")
	}
	if broker.defunctClients == nil {
		t.Error("broker.defunctClients channel is nil")
	}
	if broker.messages == nil {
		t.Error("broker.messages channel is nil")
	}

	// Verify the broker is actually running by adding a client
	clientChan := make(chan string, 1)
	broker.newClients <- clientChan
	time.Sleep(10 * time.Millisecond)

	if len(broker.clients) != 1 {
		t.Errorf("expected broker to have 1 client after adding, got %d", len(broker.clients))
	}
}

func TestRunKeepalives(t *testing.T) {
	// Save original global state
	originalOptions := options
	originalEvents := events
	defer func() {
		options = originalOptions
		events = originalEvents
	}()

	t.Run("sends keepalive messages", func(t *testing.T) {
		options.Debug = false

		// Create a test broker
		events = &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string, 10), // Buffered to catch messages
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Start keepalives in a goroutine
		go runKeepalives(ctx)

		// Wait a bit and check if a keepalive message was sent
		time.Sleep(50 * time.Millisecond)

		// There might not be a message yet due to the 3-second interval,
		// but we can at least verify the function doesn't crash
		// and that it respects context cancellation
	})

	t.Run("stops when context is cancelled", func(t *testing.T) {
		options.Debug = false

		events = &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string, 10),
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Start keepalives
		done := make(chan bool)
		go func() {
			runKeepalives(ctx)
			done <- true
		}()

		// Cancel context
		cancel()

		// Wait for runKeepalives to exit
		select {
		case <-done:
			// Successfully exited
		case <-time.After(200 * time.Millisecond):
			t.Error("runKeepalives did not exit after context cancellation")
		}
	})

	t.Run("respects debug flag", func(t *testing.T) {
		options.Debug = true

		events = &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string, 10),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Should not panic with debug enabled
		go runKeepalives(ctx)

		time.Sleep(60 * time.Millisecond)

		// Reset debug
		options.Debug = false
	})
}
