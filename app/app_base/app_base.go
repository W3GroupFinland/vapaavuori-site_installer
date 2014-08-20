package app_base

import (
	"github.com/tuomasvapaavuori/site_installer/app/modules/config"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
)

type AppBase struct {
	DataStore database.DataStore
	Config    *config.Config
	Commands  Commands
}

type Commands struct {
	HttpServer *HttpServer
}

type HttpServer struct {
	Restart *config.Command
}

func NewAppBase() *AppBase {
	a := AppBase{
		Config:    new(config.Config),
		DataStore: database.DataStore{},
		Commands:  Commands{HttpServer: &HttpServer{}},
	}

	return &a
}
