package server

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"

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

	err = json.Unmarshal(byteValue, &boxes)
	if err != nil {
		logger.Fatal(err.Error())
	}

	sortBoxes()

}

// Write json
func saveBoxFile() error {
	byteValue, err := json.Marshal(&boxes)
	if err != nil {
		return err
	}

	if _, err := os.Stat(boxFile + ".bak9"); err == nil {
		err = os.Remove(boxFile + ".bak9")
		if err != nil {
			logger.Error(err.Error())
		}
	} else if errors.Is(err, os.ErrNotExist) {
	} else {
		logger.Error(err.Error())
	}

	for i := 8; i > 0; i-- {
		s := strconv.Itoa(i)
		t := strconv.Itoa(i + 1)
		if _, err := os.Stat(boxFile + ".bak" + s); err == nil {
			os.Rename(boxFile+".bak"+s, boxFile+".bak"+t)
		} else if errors.Is(err, os.ErrNotExist) {
		} else {
			logger.Error(err.Error())
		}
	}

	if _, err := os.Stat(boxFile); err == nil {
		os.Rename(boxFile, boxFile+".bak1")
	} else if errors.Is(err, os.ErrNotExist) {
	} else {
		logger.Error(err.Error())
	}

	err = os.WriteFile(boxFile, byteValue, 0644)
	if err != nil {
		return err
	}

	return nil
}
