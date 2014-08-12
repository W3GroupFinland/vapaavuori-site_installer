package main

import (
	"github.com/tuomasvapaavuori/site_installer/app"
)

func main() {
	application := app.Init()
	application.Run()
}
