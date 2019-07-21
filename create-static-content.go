package main

import (
	"log"
	"os"
)

func createStaticContent() {

	for _, file := range AssetNames() {
		if _, err := os.Stat(config.staticFilePath + file); os.IsNotExist(err) {
			log.Printf("'%s/%s' doesn't exist, creating default file.", config.staticFilePath, file)
			RestoreAsset(config.staticFilePath, file)
		} else if config.useDefaultStatic {
			log.Printf("Default files enforced, creating default file '%s/%s'.", config.staticFilePath, file)
			RestoreAsset(config.staticFilePath, file)
		}
	}
	if _, err := os.Stat(config.dataFile); os.IsNotExist(err) {
		var file, err = os.Create(config.dataFile)
		if err != nil {
			log.Printf("Data file did not exist and could not create an empty one.")
			log.Fatal(err)
		}
		log.Printf("Created empty data file %s", config.dataFile)
		defer file.Close()
	}

}
