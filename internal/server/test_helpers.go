package server

import "go.uber.org/zap"

// Test helper functions

func initTestLogger() {
	if logger == nil {
		logger = zap.NewNop()
	}
}
