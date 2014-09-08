package app

import (
	"fmt"
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	a "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"net/http"
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

	rf := utils.NewRequireFiles()
	rf.Add(a.Base.Config.Hosts.File, "Hosts file").
		Add(a.Base.Config.Backup.Directory, "Backup directory").
		Add(a.Base.Config.Platform.Directory, "Platform directory").
		Add(a.Base.Config.SiteTemplates.Directory, "Site template directory").
		Add(a.Base.Config.SiteServerTemplates.Directory, "Site server template directory").
		Add(a.Base.Config.SiteServerTemplates.Certificates, "Site server templates certificates directory").
		Add(a.Base.Config.ServerConfigRoot.Http, "Site HTTP config root").
		Add(a.Base.Config.ServerConfigRoot.SSL, "Site SSL config root")

	if a.Base.Config.Ssl.HttpSsl {
		rf.Add(a.Base.Config.Ssl.CertFile, "SSL Certificate file").
			Add(a.Base.Config.Ssl.PrivateFile, "SSL Certificate Private file")
	}

	if !utils.UserExists(a.Base.Config.HttpUser.User) {
		log.Fatalln("Check that http user exists.")
	}

	if !utils.UserExists(a.Base.Config.DeployUser.User) {
		log.Fatalln("Check that deploy user exists.")
	}

	rf.Require()

	return &a
}

func (a *Application) RegisterControllers() {
	a.Controllers.Drush = &controllers.Drush{Base: a.Base}
	a.Controllers.Drush.Init()
	a.Controllers.Site = &controllers.Site{Drush: a.Controllers.Drush, Base: a.Base}
	a.Controllers.SiteTemplate = &controllers.SiteTemplate{Base: a.Base}
	a.Controllers.HostMasterDB = &controllers.HostMasterDB{Base: a.Base}

	a.Controllers.System = &controllers.System{
		HostMaster: a.Controllers.HostMasterDB,
		Site:       a.Controllers.Site,
	}
}

func (a *Application) Run() {
	// Get command line flags.
	a.GetAppCommandArgs()

	defer a.Base.DataStore.DB.Close()

	// Read web templates.
	templatePath, err := utils.GetFileFullPath("web/templates")
	if err != nil {
		log.Fatalln("Application web template folder doesn't exist.")
	}

	err = a.Base.Templates.CustomDelims().ReadDir(templatePath)
	if err != nil {
		log.Println(err)
	}

	// Register controllers.
	a.RegisterControllers()
	a.RegisterWebControllers()

	// Register routes.
	a.RegisterRoutes()

	// Register file servers.
	a.RegisterPublicFileServer()
	a.RegisterWebAppServer()

	a.Controllers.Drush.Which()

	a.Base.AppKeys.SetSecret("client.secret", "something-wery-secret")
	a.Base.InitSessions("something-wery-secret")

	log.Printf("Running on port %d\n", *a.Base.Flags.Port)
	addr := fmt.Sprintf("%v:%d", a.Base.Config.Host.Name, *a.Base.Flags.Port)

	// Check for https settings.
	if a.Base.Config.Ssl.HttpSsl {
		err := http.ListenAndServeTLS(addr, a.Base.Config.Ssl.CertFile, a.Base.Config.Ssl.PrivateFile, nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}
