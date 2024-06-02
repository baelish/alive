package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var boxFile string

func createDataFiles() {
	if options.Debug == true {
		log.Print("Creating data files")
	}
	boxFile = filepath.Clean(options.DataPath + "/boxes.json")
	if _, err := os.Stat(options.DataPath); os.IsNotExist(err) {
		err := os.Mkdir(options.DataPath, 0755)
		if err != nil {
			log.Printf("Data directory didn't exist and couldn't create it (%s)", options.DataPath)
		}
	}

	if _, err := os.Stat(boxFile); os.IsNotExist(err) {
		var file, err = os.Create(boxFile)
		if err != nil {
			log.Printf("Data file did not exist and could not create an empty one.")
			log.Fatal(err)
		}

		err = os.WriteFile(boxFile, []byte(emptyDataFile), 0644)
		if err != nil {
			log.Printf("Could not add base content to file %s", boxFile)
			log.Fatal(err)
		}

		log.Printf("Created empty data file %s", boxFile)
		defer func() {
			err = file.Close()
		}()
		if err != nil {
			log.Print(err)
		}
	}
}

// Loads Json from a file and returns Boxes sorted by size (Largest first)
func getBoxesFromDataFile() {
	if options.Debug == true {
		log.Print("Getting boxes from data file")
	}
	byteValue, err := os.ReadFile(boxFile)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(byteValue, &boxes)
	if err != nil {
		log.Fatal(err)
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
			log.Print(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
	} else {
		log.Print(err)
	}

	for i := 8; i > 0; i-- {
		s := strconv.Itoa(i)
		t := strconv.Itoa(i + 1)
		if _, err := os.Stat(boxFile + ".bak" + s); err == nil {
			os.Rename(boxFile+".bak"+s, boxFile+".bak"+t)
		} else if errors.Is(err, os.ErrNotExist) {
		} else {
			log.Print(err)
		}
	}

	if _, err := os.Stat(boxFile); err == nil {
		os.Rename(boxFile, boxFile+".bak1")
	} else if errors.Is(err, os.ErrNotExist) {
	} else {
		log.Print(err)
	}

	err = os.WriteFile(boxFile, byteValue, 0644)
	if err != nil {
		return err
	}

	return nil
}
