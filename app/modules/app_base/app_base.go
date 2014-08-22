package app_base

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/tuomasvapaavuori/site_installer/app/modules/config"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	t "github.com/tuomasvapaavuori/site_installer/app/modules/templates"
)

type AppBase struct {
	Templates      t.Templates
	Routes         map[string]*Route
	WebControllers Controllers
	AppKeys        AppKeys
	Config         *config.Config
	Sessions       *sessions.CookieStore
	JSON           JSON
	Http           Http
	DataStore      database.DataStore
	Commands       Commands
}

type Commands struct {
	HttpServer *HttpServer
}

type HttpServer struct {
	Restart *config.Command
}

type AppKeys struct {
	Secrets map[string]string
}

func NewAppBase() *AppBase {
	a := AppBase{
		Config:    new(config.Config),
		DataStore: database.DataStore{},
		Commands:  Commands{HttpServer: &HttpServer{}},
	}

	return &a
}

func (a *AppKeys) SetSecret(name string, value string) {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}

	hasher := sha1.New()
	hasher.Write([]byte(value))

	a.Secrets[name] = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (a *AppKeys) GetSecret(name string) (string, error) {
	if _, exists := a.Secrets[name]; !exists {
		return "", errors.New(fmt.Sprintf("Secret %v not set.", name))
	}

	return a.Secrets[name], nil
}
