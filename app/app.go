package app

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
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
		HostMasterDB *controllers.HostMasterDB
	}
}

func Init(config []byte) *Application {
	// Create new application struct with application base.
	a := Application{Base: a.NewAppBase()}
	// Read application configuration from file.
	a.Base.Config.Read(config)
	// Open database connection.
	_, err := a.Base.DataStore.OpenConn(
		a.Base.Config.Mysql.User,
		a.Base.Config.Mysql.Password,
		a.Base.Config.Mysql.Protocol,
		a.Base.Config.Mysql.Host,
		a.Base.Config.Mysql.Port,
		a.Base.Config.Mysql.DbName)

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := a.Base.Config.HttpServer.GetRestartCmd()
	if err == nil {
		log.Println(err)
	}

	a.Base.Commands.HttpServer.Restart = cmd

	if !utils.FileExists(a.Base.Config.Hosts.Directory) {
		log.Fatalln("Hosts file doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.Backup.Directory) {
		log.Fatalln("Backup directory doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.Platform.Directory) {
		log.Fatalln("Platform directory doesn't exist. Please correct it before continuing.")
	}

	return &a
}

func (a *Application) RegisterControllers() {
	a.Controllers.Drush = &controllers.Drush{Base: a.Base}
	a.Controllers.Drush.Init()
	a.Controllers.Site = &controllers.Site{Drush: a.Controllers.Drush, Base: a.Base}
	a.Controllers.SiteTemplate = &controllers.SiteTemplate{Base: a.Base}
	a.Controllers.System = &controllers.System{a.Controllers.Site}
	a.Controllers.HostMasterDB = &controllers.HostMasterDB{Base: a.Base}
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

	tmpl.InstallInfo.DomainInfo = &models.Domain{Host: "127.0.0.1", DomainName: "local.tivia-drupal1.fi"}
	tmpl.HttpServer.DomainInfo = &models.Domain{Host: "127.0.0.1", DomainName: "local.tivia-drupal1.fi"}
	tmpl.SSLServer.DomainInfo = &models.Domain{Host: "127.0.0.1", DomainName: "local.ssl-tivia-drupal1.fi"}

	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(tmpl)

	_, err = a.Controllers.Site.Create(tmpl)
	if err != nil {
		log.Println(err)
		tmpl.RollBack.Execute()
		return
	}

	err = a.Controllers.SiteTemplate.WriteApacheConfig(tmpl)
	if err != nil {
		log.Println(err)
		tmpl.RollBack.Execute()
		return
	}

	domains := a.Controllers.Site.GetSiteTemplateDomains(tmpl)
	err = a.Controllers.Site.AddToHosts(tmpl, domains)
	if err != nil {
		log.Println(err)
		tmpl.RollBack.Execute()
		return
	}

	a.Controllers.Site.CreateDomainSymlinks(tmpl, domains)

	err = a.Controllers.System.HttpServerRestart()
	if err != nil {
		log.Println(err)
		return
	}

	tmpl.RollBack.DeleteBackupFiles()
}
