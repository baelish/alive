package server

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/baelish/alive/api"
	"go.uber.org/zap"
)

var boxFile string

func createDataFiles() {
	logger.Debug("Creating data files")
	boxFile = filepath.Clean(options.DataPath + "/boxes.json")
	if _, err := os.Stat(options.DataPath); os.IsNotExist(err) {
		err := os.Mkdir(options.DataPath, 0755)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	if _, err := os.Stat(boxFile); os.IsNotExist(err) {
		var file, err = os.Create(boxFile)
		if err != nil {
			logger.Fatal(err.Error())
		}

		err = os.WriteFile(boxFile, []byte(emptyDataFile), 0644)
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Created empty data file", zap.String("file", boxFile))
		defer func() {
			err = file.Close()
		}()
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

// Loads Json from a file and returns Boxes sorted by size (Largest first)
func getBoxesFromDataFile() {
	if options.Debug {
		logger.Info("Getting boxes from data file")
	}
	byteValue, err := os.ReadFile(boxFile)

	if err != nil {
		logger.Fatal(err.Error())
	}

	var loadedBoxes []api.Box
	err = json.Unmarshal(byteValue, &loadedBoxes)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Load boxes into the store (thread-safe)
	boxStore.mu.Lock()
	boxStore.boxes = loadedBoxes
	boxStore.sortUnsafe()
	boxStore.mu.Unlock()
}

// saveBoxFile writes all boxes to the data file, keeping up to nine backups.
// The new data is written to a temporary file and renamed over the data file,
// so a crash or kill part way through never leaves it missing or truncated.
func saveBoxFile() error {
	// Get all boxes from store (thread-safe)
	boxes := boxStore.GetAll()
	byteValue, err := json.Marshal(&boxes)
	if err != nil {
		return err
	}

	tmpFile := boxFile + ".tmp"
	if err := writeFileSynced(tmpFile, byteValue); err != nil {
		os.Remove(tmpFile)
		return err
	}

	if err := os.Remove(boxFile + ".bak9"); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error(err.Error())
	}

	for i := 8; i > 0; i-- {
		s := strconv.Itoa(i)
		t := strconv.Itoa(i + 1)
		if err := os.Rename(boxFile+".bak"+s, boxFile+".bak"+t); err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Error(err.Error())
		}
	}

	// Keep the current data file as the newest backup. A hard link leaves the
	// data file in place until the rename below replaces it in one step.
	if err := os.Link(boxFile, boxFile+".bak1"); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Warn("could not link backup, falling back to rename", zap.Error(err))
		if err := os.Rename(boxFile, boxFile+".bak1"); err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Error(err.Error())
		}
	}

	return os.Rename(tmpFile, boxFile)
}

// writeFileSynced writes data to name and flushes it to disk before closing.
func writeFileSynced(name string, data []byte) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
