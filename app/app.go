package app

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	//"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"log"
)

type Application struct {
	Base      *a.AppBase
	Arguments map[string]string
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
	a.ParseCommandLineArgs()
	/*newDB := database.NewDatabase(&a.Base.DataStore)
	db, err := newDB.Randomize().
		SetHosts([]string{"localhost", "127.0.0.1"}).
		SetUserPrivileges([]string{"ALL"}, true).
		CreateDatabase()

	if err != nil {
		log.Println(err, db)
	}*/

	val, err := a.GetCommandArg("--template")
	if err != nil {
		log.Println(err)
	}
	log.Println(val)

	defer a.Base.DataStore.DB.Close()
}
