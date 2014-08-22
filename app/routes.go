package app

import (
	"github.com/tuomasvapaavuori/site_installer/app/controllers/web"
	base "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"net/http"
)

func (a *Application) RegisterWebControllers() {
	app := a.Base

	user := &controllers.User{app}
	app.WebControllers.Set(&base.Controller{user})
	app.WebControllers.Set(&base.Controller{&controllers.Login{user}})
}

func (a *Application) RegisterRoutes() {
	app := a.Base

	app.AddRoute(&base.Route{
		Path: "/user/login",
		Handlers: []base.RouteHandler{
			base.RouteHandler{
				Method:  "GET",
				Handler: app.WebControllers.Get("app.web.controllers.login").Handler("Page")},
			base.RouteHandler{
				Method:  "POST",
				Handler: app.WebControllers.Get("app.web.controllers.login").Handler("LoginPostHandler")},
		}})

	app.AddRoute(&base.Route{
		Path: "/user/logout",
		Handlers: []base.RouteHandler{
			base.RouteHandler{
				Method:  "GET",
				Handler: app.WebControllers.Get("app.web.controllers.login").Handler("Logout")}}})

	// Register application routes
	app.RegisterRoutes()
}

func (a *Application) RegisterFileServer(dir *string) {
	// File server
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/files/", http.StripPrefix("/files/", fileHandler))
}
