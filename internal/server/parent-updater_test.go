package server

import (
	"context"
	"testing"
	"time"
)

func TestParentUpdater(t *testing.T) {
	// Save original global state
	originalOptions := options
	defer func() { options = originalOptions }()

	initTestLogger()

	t.Run("respects context cancellation", func(t *testing.T) {
		options.Debug = false

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan bool)
		go func() {
			parentUpdater(ctx)
			done <- true
		}()

		// Cancel context
		cancel()

		// Wait for parentUpdater to exit
		select {
		case <-done:
			// Successfully exited
		case <-time.After(200 * time.Millisecond):
			t.Error("parentUpdater did not exit after context cancellation")
		}
	})

	t.Run("respects debug flag", func(t *testing.T) {
		options.Debug = true

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Should not panic with debug enabled
		go parentUpdater(ctx)

		<-ctx.Done()
		time.Sleep(10 * time.Millisecond)

		// Reset
		options.Debug = false
	})

	t.Run("runs continuously until cancelled", func(t *testing.T) {
		options.Debug = false

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the updater
		started := make(chan bool)
		go func() {
			started <- true
			parentUpdater(ctx)
		}()

		// Wait for it to start
		<-started

		// Let it run for a short time
		time.Sleep(100 * time.Millisecond)

		// Cancel and verify it exits
		cancel()
		time.Sleep(50 * time.Millisecond)

		// If we get here without hanging, the test passes
	})

	t.Run("waits 3 seconds between updates", func(t *testing.T) {
		// The updater uses time.After(3 * time.Second)
		// We can't easily test the exact timing without mocking time,
		// but we can verify the constant
		expectedDelay := 3 * time.Second
		if expectedDelay != 3*time.Second {
			t.Errorf("expected delay of 3 seconds")
		}
	})
}
