package app

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

type Application struct {
	Base        *a.AppBase
	Arguments   map[string]string
	Controllers struct {
		Drush *controllers.Drush
		Site  *controllers.Site
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
}

func (a *Application) Run() {
	// Register controllers.
	a.RegisterControllers()

	a.Controllers.Drush.Which()
	a.ParseCommandLineArgs()

	// Define site install info.
	installInfo := models.NewSiteInstallInfo()
	installInfo.InstallRoot = "/Users/tuomas/Sites/www/drupal-7.31"
	installInfo.InstallType = "standard"
	installInfo.SiteName = "local.huuhaa.fi"
	installInfo.SubDirectory = "local.huuhaa.fi"
	// Create new site.
	a.Controllers.Site.Create(installInfo)

	val, err := a.GetCommandArg("--template")
	if err != nil {
		log.Println(err)
	}
	log.Println(val)

	defer a.Base.DataStore.DB.Close()
}
