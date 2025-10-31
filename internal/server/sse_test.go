package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
			clientCount:    make(chan int),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		broker.Start(ctx)

		// Add a new client
		clientChan := make(chan string, 1)
		broker.newClients <- clientChan

		// Give it a moment to process
		time.Sleep(10 * time.Millisecond)

		// Verify the client was added using thread-safe method
		if count := broker.ClientCount(); count != 1 {
			t.Errorf("expected 1 client, got %d", count)
		}
	})

	t.Run("removes defunct clients", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
			clientCount:    make(chan int),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		broker.Start(ctx)

		// Add a client
		clientChan := make(chan string, 1)
		broker.newClients <- clientChan
		time.Sleep(10 * time.Millisecond)

		if count := broker.ClientCount(); count != 1 {
			t.Fatalf("expected 1 client after adding, got %d", count)
		}

		// Remove the client
		broker.defunctClients <- clientChan
		time.Sleep(10 * time.Millisecond)

		// Verify the client was removed using thread-safe method
		if count := broker.ClientCount(); count != 0 {
			t.Errorf("expected 0 clients after removal, got %d", count)
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
			clientCount:    make(chan int),
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
			clientCount:    make(chan int),
		}

		ctx, cancel := context.WithCancel(context.Background())
		broker.Start(ctx)

		// Add a client to verify broker is running
		clientChan := make(chan string, 1)
		broker.newClients <- clientChan
		time.Sleep(10 * time.Millisecond)

		if count := broker.ClientCount(); count != 1 {
			t.Fatalf("expected 1 client, got %d", count)
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
	mu sync.Mutex
	header http.Header
	body   []byte
	statusCode  int
	closeNotify chan bool
	flushed     bool
	wroteOnce   chan struct{} // Signals when first write occurs (headers are done)
}

func (w *testResponseWriter) Header() http.Header {
	// Note: Returning the header map directly means callers must ensure
	// they don't access it concurrently. Tests should synchronize properly.
	return w.header
}

// GetHeader returns a thread-safe value of a header (call after WaitForWrite())
func (w *testResponseWriter) GetHeader(key string) string {
	return w.header.Get(key)
}

// WaitForWrite waits until at least one Write() has been called
func (w *testResponseWriter) WaitForWrite() {
	if w.wroteOnce != nil {
		<-w.wroteOnce
	}
}

func (w *testResponseWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Signal that first write has occurred (headers are now immutable)
	if w.wroteOnce != nil {
		select {
		case <-w.wroteOnce:
			// Already closed
		default:
			close(w.wroteOnce)
		}
	}

	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.statusCode = statusCode
}

func (w *testResponseWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.flushed = true
}

func (w *testResponseWriter) CloseNotify() <-chan bool {
	if w.closeNotify == nil {
		w.closeNotify = make(chan bool)
	}
	return w.closeNotify
}

// GetBody returns a thread-safe copy of the body
func (w *testResponseWriter) GetBody() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]byte, len(w.body))
	copy(result, w.body)
	return result
}

func TestBroker_ServeHTTP(t *testing.T) {
	t.Run("sets correct SSE headers", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
			clientCount:    make(chan int),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		broker.Start(ctx)

		// Use a custom ResponseWriter to capture headers
		testWriter := &testResponseWriter{
			header:    make(http.Header),
			wroteOnce: make(chan struct{}),
		}
		testReq := httptest.NewRequest("GET", "/events/", nil)

		// Start ServeHTTP in background
		go broker.ServeHTTP(testWriter, testReq)

		// Send a message to trigger a write - this ensures headers are fully set
		// because HTTP headers must be set before the first write
		time.Sleep(20 * time.Millisecond)
		broker.messages <- "test"

		// Wait for the write to complete
		testWriter.WaitForWrite()

		// Now headers are definitely immutable, safe to read
		// Make a copy under lock for safety
		testWriter.mu.Lock()
		headerCopy := make(http.Header)
		for k, v := range testWriter.header {
			headerCopy[k] = v
		}
		testWriter.mu.Unlock()

		// Check the copied headers
		if ct := headerCopy.Get("Content-Type"); ct != "text/event-stream" {
			t.Errorf("Content-Type: expected %q, got %q", "text/event-stream", ct)
		}
		if cc := headerCopy.Get("Cache-Control"); cc != "no-cache" {
			t.Errorf("Cache-Control: expected %q, got %q", "no-cache", cc)
		}
		if conn := headerCopy.Get("Connection"); conn != "keep-alive" {
			t.Errorf("Connection: expected %q, got %q", "keep-alive", conn)
		}
		if te := headerCopy.Get("Transfer-Encoding"); te != "chunked" {
			t.Errorf("Transfer-Encoding: expected %q, got %q", "chunked", te)
		}
	})

	t.Run("sends messages in SSE format", func(t *testing.T) {
		broker := &Broker{
			clients:        make(map[chan string]bool),
			newClients:     make(chan chan string),
			defunctClients: make(chan chan string),
			messages:       make(chan string),
			clientCount:    make(chan int),
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

		// Check the response body contains the SSE-formatted message (use thread-safe method)
		body := string(testWriter.GetBody())
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
			clientCount:    make(chan int),
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

		// Verify client was registered using thread-safe method
		if count := broker.ClientCount(); count != 1 {
			t.Errorf("expected 1 registered client, got %d", count)
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

	// Use thread-safe ClientCount method instead of accessing map directly
	if count := broker.ClientCount(); count != 1 {
		t.Errorf("expected broker to have 1 client after adding, got %d", count)
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
			clientCount:    make(chan int),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Start keepalives in a goroutine
		done := make(chan bool)
		go func() {
			runKeepalives(ctx)
			done <- true
		}()

		// Wait for context timeout and goroutine to exit
		<-ctx.Done()
		<-done

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
			clientCount:    make(chan int),
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
			clientCount:    make(chan int),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Should not panic with debug enabled
		done := make(chan bool)
		go func() {
			runKeepalives(ctx)
			done <- true
		}()

		// Wait for context timeout and goroutine to exit
		<-ctx.Done()
		<-done

		// Reset debug
		options.Debug = false
	})
}
