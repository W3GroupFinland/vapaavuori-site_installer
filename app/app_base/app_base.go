package app_base

import (
	"github.com/tuomasvapaavuori/site_installer/app/modules/config"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
)

type AppBase struct {
	DataStore database.DataStore
	Config    *config.Config
}

func NewAppBase() *AppBase {
	a := AppBase{
		Config:    new(config.Config),
		DataStore: database.DataStore{},
	}

	return &a
}
