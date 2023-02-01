package main

import (
	"log"
	"os"
)

func createStaticContent() {
	for _, file := range AssetNames() {
		if _, err := os.Stat(options.StaticPath + "/" + file); os.IsNotExist(err) {
			log.Printf("'%s/%s' doesn't exist, creating default file.", options.StaticPath, file)
			err = RestoreAsset(options.StaticPath, file)
			if err != nil {
				log.Print(err)
			}
		} else if options.DefaultStatic {
			log.Printf("Default files enforced, creating default file '%s/%s'.", options.StaticPath, file)
			err = RestoreAsset(options.StaticPath, file)
			if err != nil {
				log.Print(err)
			}
		}
	}

	if _, err := os.Stat(options.DataFile); os.IsNotExist(err) {
		var file, err = os.Create(options.DataFile)

		if err != nil {
			log.Printf("Data file did not exist and could not create an empty one.")
			log.Fatal(err)
		}

		log.Printf("Created empty data file %s", options.DataFile)
		defer func() {
			err = file.Close()
		}()
		if err != nil {
			log.Print(err)
		}
	}

}
