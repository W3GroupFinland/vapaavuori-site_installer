package app

import (
	"github.com/tuomasvapaavuori/site_installer/app/controllers/web"
	base "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"net/http"
)

func (a *Application) RegisterWebControllers() {
	app := a.Base

	user := &controllers.User{app}
	// User controller
	app.WebControllers.Set(&base.Controller{user})
	// Login controller
	//login.User = user
	app.WebControllers.Set(&base.Controller{&controllers.Login{User: user}})
	// Application page controller
	app.WebControllers.Set(&base.Controller{&controllers.ApplicationPage{user}})
	// Hostmaster websocket controller
	app.WebControllers.Set(&base.Controller{&controllers.HostmasterWS{
		User:         user,
		Base:         app,
		System:       a.Controllers.System,
		HostMasterDB: a.Controllers.HostMasterDB,
		Site:         a.Controllers.Site,
		SiteTemplate: a.Controllers.SiteTemplate,
	}})
}

func (a *Application) RegisterRoutes() {
	app := a.Base

	app.AddRoute(&base.Route{
		Path: "/",
		Handlers: []base.RouteHandler{
			base.RouteHandler{
				Method:  "GET",
				Handler: app.WebControllers.Get("app.controllers.web.login").Handler("Page")},
			base.RouteHandler{
				Method:  "POST",
				Handler: app.WebControllers.Get("app.controllers.web.login").Handler("LoginPostHandler")},
		}})

	app.AddRoute(&base.Route{
		Path: "/logout",
		Handlers: []base.RouteHandler{
			base.RouteHandler{
				Method:  "GET",
				Handler: app.WebControllers.Get("app.controllers.web.login").Handler("Logout")}}})

	app.AddRoute(&base.Route{
		Path: "/app/ws",
		Handlers: []base.RouteHandler{
			base.RouteHandler{
				Method:  "GET",
				Handler: app.WebControllers.Get("app.controllers.web.hostmaster.ws").Handler("Messager")}},
		Acl: []base.HttpAclHandlerFunc{app.WebControllers.Get("app.controllers.web.user").AclHandler("LoggedInAcl")}})

	// Register application routes
	app.RegisterRoutes()
}

func (a *Application) RegisterPublicFileServer() {
	fullPath, err := utils.GetFileFullPath("web/files/public")
	if err != nil {
		log.Fatalln(err)
	}

	// File server
	fs := base.NewFileServer(fullPath, "/public/")

	http.Handle("/public/", fs)
}

func (a *Application) RegisterWebAppServer() {
	app := a.Base

	// If production mode, serve files from dist.
	ap := "web/files/webapp/dist/"
	if *a.Base.Flags.DevMode {
		// If development mode serve from webapp.
		ap = "web/files/webapp/app/"
		log.Println("Serving application files in development mode.")
	}

	fullPath, err := utils.GetFileFullPath(ap)
	if err != nil {
		log.Fatalln(err)
	}

	// Application server
	as := base.NewFileServer(fullPath, "/app/").
		AddAcl(app.WebControllers.Get("app.controllers.web.user").AclHandler("LoggedInAcl"))

	http.Handle("/app/", as)
}
