package app

import (
	"flag"
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

	if !utils.FileExists(a.Base.Config.Hosts.Directory) {
		log.Fatalln("Hosts file doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.Backup.Directory) {
		log.Fatalln("Backup directory doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.Platform.Directory) {
		log.Fatalln("Platform directory doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.SiteTemplates.Directory) {
		log.Fatalln("Site template directory doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.SiteServerTemplates.Directory) {
		log.Fatalln("Site server template directory doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.SiteServerTemplates.Certificates) {
		log.Fatalln("Site server templates directory doesn't exist. Please correct it before continuing.")
	}

	if !utils.FileExists(a.Base.Config.ServerConfigRoot.Directory) {
		log.Fatalln("Site server config root directory doesn't exist. Please correct it before continuing.")
	}

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
	defer a.Base.DataStore.DB.Close()

	// Read web templates.
	err := a.Base.Templates.CustomDelims().ReadDir("web/templates")
	if err != nil {
		log.Println(err)
	}

	// Command line flags
	port := flag.Int("port", a.Base.Config.Host.Port, "port to serve on")
	dir := flag.String("directory", "web/files", "directory of web files")
	flag.Parse()

	// Register controllers.
	a.RegisterControllers()
	a.RegisterWebControllers()
	a.RegisterRoutes()
	a.RegisterFileServer(dir)

	a.Controllers.Drush.Which()
	a.ParseCommandLineArgs()

	a.Base.AppKeys.SetSecret("client.secret", "something-wery-secret")
	a.Base.InitSessions("something-wery-secret")

	log.Printf("Running on port %d\n", *port)
	addr := fmt.Sprintf("%v:%d", a.Base.Config.Host.Name, *port)
	err = http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
