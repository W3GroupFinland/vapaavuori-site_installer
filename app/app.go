package app

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	"log"
)

type Application struct {
	Base        *a.AppBase
	Arguments   map[string]string
	Controllers struct {
		Drush        *controllers.Drush
		Site         *controllers.Site
		SiteTemplate *controllers.SiteTemplate
		System       *controllers.System
	}
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

func (a *Application) RegisterControllers() {
	a.Controllers.Drush = &controllers.Drush{Base: a.Base}
	a.Controllers.Drush.Init()
	a.Controllers.Site = &controllers.Site{Drush: a.Controllers.Drush, Base: a.Base}
	a.Controllers.SiteTemplate = &controllers.SiteTemplate{Base: a.Base}
	a.Controllers.System = &controllers.System{Base: a.Base}
}

func (a *Application) Run() {
	defer a.Base.DataStore.DB.Close()
	// Register controllers.
	a.RegisterControllers()

	a.Controllers.Drush.Which()
	a.ParseCommandLineArgs()

	templFile, err := a.GetCommandArg("--template")
	if err != nil {
		log.Println(err)
		return
	}

	tmpl, err := a.Controllers.SiteTemplate.ReadTemplate(templFile)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = a.Controllers.Site.Create(tmpl)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.Controllers.SiteTemplate.WriteApacheConfig(tmpl)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.Controllers.Site.AddToHosts(tmpl)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.Controllers.System.ApacheRestart()
	if err != nil {
		log.Println(err)
		return
	}
}
