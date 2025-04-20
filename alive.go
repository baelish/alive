package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const timeFormat = "2006-01-02T15:04:05.000Z07:00"

var events *Broker

func main() {
	logger = zap.Must(zap.NewProduction())

	processOptions()

	if os.Getenv("DEV") != "" {
		logger.Sync()
		logger = zap.Must(zap.NewDevelopment())
	} else if options.Debug {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger.Sync()
		logger = zap.Must(cfg.Build())
	}

	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	if options.Demo {
		tempDir, err := os.MkdirTemp(os.TempDir(), "alive-*.tmp")
		if err != nil {
			logger.Panic("Unable to create a temporary directory", zap.String("dir", tempDir))
		}
		defer os.RemoveAll(tempDir)

		logger.Info("Running demo using temporary files", zap.String("dir", tempDir))

		options.DataPath = filepath.Clean(fmt.Sprintf("%s/data", tempDir))
		options.StaticPath = filepath.Clean(fmt.Sprintf("%s/static", tempDir))
	}

	if options.DataPath == "" {
		options.DataPath = filepath.Clean(fmt.Sprintf("%s/.alive/data", os.Getenv("HOME")))
	}

	if options.StaticPath == "" {
		options.StaticPath = filepath.Clean(fmt.Sprintf("%s/.alive/static", os.Getenv("HOME")))
	}

	logger.Debug("options requested", logStructDetails(options)...)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		logger.Info("Running cleanup")
		cancel()
	}()

	createStaticContent()
	createDataFiles()
	getBoxesFromDataFile()

	events = runSSE(ctx)

	go runDashboard(ctx)
	go runAPI(ctx)

	go runKeepalives(ctx)
	go maintainBoxes(ctx)

	if options.ParentUrl != "" {
		go parentUpdater(ctx)
	}

	if options.Demo {
		go runDemo(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(1 * time.Second)):
		}
	}
}
