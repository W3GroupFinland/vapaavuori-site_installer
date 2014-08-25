package app

import (
	"github.com/tuomasvapaavuori/site_installer/app/controllers/web"
	base "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"net/http"
)

func (a *Application) RegisterWebControllers() {
	app := a.Base

	user := &controllers.User{app}
	// User controller
	app.WebControllers.Set(&base.Controller{user})
	// Login controller
	app.WebControllers.Set(&base.Controller{&controllers.Login{user}})
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
		Path: "/app",
		Handlers: []base.RouteHandler{
			base.RouteHandler{
				Method:  "GET",
				Handler: app.WebControllers.Get("app.controllers.web.application.page").Handler("ApplicationPageHandler")}},
		Acl: []base.HttpAclHandlerFunc{app.WebControllers.Get("app.controllers.web.user").AclHandler("LoggedInAcl")}})

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

func (a *Application) RegisterFileServer(dir *string) {
	// File server
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/files/", http.StripPrefix("/files/", fileHandler))
}
