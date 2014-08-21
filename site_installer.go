package main

import (
	"github.com/tuomasvapaavuori/site_installer/app"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
)

func main() {
	// Abstracted config reading to application root level.
	// Makes it more easy to test application with string config..
	config, err := utils.ReadFileContents("config/config.gcfg")
	if err != nil {
		log.Fatalln(err)
	}

	application := app.Init(config)
	application.Run()
}
