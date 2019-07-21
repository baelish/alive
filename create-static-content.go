package main

import (
	"log"
	"os"
)

func createDataFile(dataFile string) {
	var _, err = os.Stat(dataFile)
	if os.IsNotExist(err) {
		var file, err = os.Create(dataFile)
		if err != nil {
			log.Printf("Data file did not exist and could not create an empty one.")
			log.Fatal(err)
		}
		defer file.Close()
		log.Printf("Created empty data file %s", dataFile)
	}
}

func createStaticContent(path string) {
	for _, file := range AssetNames() {
		if _, err := os.Stat(path + "/" + file); os.IsNotExist(err) {
			log.Printf("'%s/%s' doesn't exist, creating default file.", path, file)
			RestoreAsset(path, file)
		} else if config.useDefaultStatic {
			log.Printf("Default files enforced, creating default file '%s/%s'.", path, file)
			RestoreAsset(path, file)
		}
	}
}
