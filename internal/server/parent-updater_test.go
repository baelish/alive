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
		done := make(chan bool)
		go func() {
			parentUpdater(ctx)
			done <- true
		}()

		<-ctx.Done()
		// Wait for goroutine to fully exit before modifying options
		<-done

		// Reset
		options.Debug = false
	})

	t.Run("runs continuously until cancelled", func(t *testing.T) {
		options.Debug = false

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the updater
		done := make(chan bool)
		go func() {
			parentUpdater(ctx)
			done <- true
		}()

		// Let it run for a short time
		time.Sleep(100 * time.Millisecond)

		// Cancel and verify it exits
		cancel()

		// Wait for goroutine to exit
		select {
		case <-done:
			// Successfully exited
		case <-time.After(100 * time.Millisecond):
			t.Error("parentUpdater did not exit after context cancellation")
		}
	})
}
