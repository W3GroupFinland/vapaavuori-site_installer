package app

import (
	"flag"
)

func (a *Application) GetAppCommandArgs() {
	// Command line flags
	a.Base.Flags.Port = flag.Int("port", a.Base.Config.Host.Port, "Port to serve on.")
	a.Base.Flags.DevMode = flag.Bool("dev-mode", false, "Defines application development mode.")

	flag.Parse()
}
