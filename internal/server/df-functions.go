package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

	// A missing data file with backups present means a save was interrupted,
	// leave it missing so the boxes are loaded from the newest backup rather
	// than starting with an empty dashboard.
	if _, err := os.Stat(boxFile); os.IsNotExist(err) && !backupExists() {
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

// maxBackups is the number of boxes.json.bakN files kept by saveBoxFile.
const maxBackups = 9

// backupFile returns the name of backup n (1 is the newest).
func backupFile(n int) string {
	return boxFile + ".bak" + strconv.Itoa(n)
}

func backupExists() bool {
	for n := 1; n <= maxBackups; n++ {
		if _, err := os.Stat(backupFile(n)); err == nil {
			return true
		}
	}
	return false
}

// readBoxFile reads and parses one data file.
func readBoxFile(name string) ([]api.Box, error) {
	byteValue, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var boxes []api.Box
	if err := json.Unmarshal(byteValue, &boxes); err != nil {
		return nil, err
	}
	return boxes, nil
}

// loadNewestBoxFile returns the boxes from the data file or, if it is missing
// or cannot be parsed, from the newest backup that can. It returns the name of
// the file used, and an error listing every file tried if none could be used.
func loadNewestBoxFile() ([]api.Box, string, error) {
	var problems []error
	for n := 0; n <= maxBackups; n++ {
		name := boxFile
		if n > 0 {
			name = backupFile(n)
		}
		boxes, err := readBoxFile(name)
		if err == nil {
			if n > 0 {
				logger.Warn("data file could not be loaded, using the newest usable backup",
					zap.String("backup", name), zap.Errors("skipped", problems))
				preserveCorruptBoxFile()
			}
			return boxes, name, nil
		}
		if n > 0 && errors.Is(err, os.ErrNotExist) {
			continue
		}
		problems = append(problems, fmt.Errorf("%s: %w", filepath.Base(name), err))
	}
	return nil, "", errors.Join(problems...)
}

// preserveCorruptBoxFile moves an unusable data file aside so that it is kept
// for inspection rather than rotated into the backups by the next save.
func preserveCorruptBoxFile() {
	if _, err := os.Stat(boxFile); err != nil {
		return
	}
	corrupt := boxFile + ".corrupt-" + time.Now().Format("20060102-150405")
	if err := os.Rename(boxFile, corrupt); err != nil {
		logger.Error("could not move the unusable data file aside", zap.Error(err))
		return
	}
	logger.Warn("moved the unusable data file aside", zap.String("file", corrupt))
}

// getBoxesFromDataFile loads the boxes into the store, falling back to the
// newest usable backup if the data file is missing or corrupt.
func getBoxesFromDataFile() {
	if options.Debug {
		logger.Info("Getting boxes from data file")
	}

	loadedBoxes, source, err := loadNewestBoxFile()
	if err != nil {
		logger.Fatal("no usable data file or backup, fix or remove them to start", zap.Error(err))
	}
	logger.Info("loaded boxes", zap.String("file", source), zap.Int("count", len(loadedBoxes)))

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

	if err := os.Remove(backupFile(maxBackups)); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error(err.Error())
	}

	for i := maxBackups - 1; i > 0; i-- {
		if err := os.Rename(backupFile(i), backupFile(i+1)); err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Error(err.Error())
		}
	}

	// Keep the current data file as the newest backup. A hard link leaves the
	// data file in place until the rename below replaces it in one step.
	if err := os.Link(boxFile, backupFile(1)); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Warn("could not link backup, falling back to rename", zap.Error(err))
		if err := os.Rename(boxFile, backupFile(1)); err != nil && !errors.Is(err, os.ErrNotExist) {
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
