package app

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/modules/user"
)

type Application struct {
	Base *a.AppBase
}

func Init() *Application {
	// Create new application struct with application base.
	a := Application{Base: a.NewAppBase()}
	// Read application configuration from file.
	a.Base.Config.Read("config/config.gcfg")
	// Open database connection.
	a.Base.DataStore.OpenConn(
		a.Base.Config.Mysql.User,
		a.Base.Config.Mysql.Password,
		a.Base.Config.Mysql.Protocol,
		a.Base.Config.Mysql.Host,
		a.Base.Config.Mysql.Port,
		a.Base.Config.Mysql.DbName)
	return &a
}

func (a *Application) Run() {
	//db.CreateDatabase("testitieto")
	u := user.User{Username: "hostmaster2", Password: "fastword"}
	a.Base.DataStore.CreateUserOnHosts(&u, []string{"localhost", "127.0.0.1"})
	a.Base.DataStore.GrantUserPrivilegesOnHosts(&u, "testitieto2", []string{"localhost", "127.0.0.1"}, []string{"ALL"}, true)
	defer a.Base.DataStore.DB.Close()
}
