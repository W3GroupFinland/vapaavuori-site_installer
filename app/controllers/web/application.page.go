package controllers

import (
	"net/http"
)

type ApplicationPage struct {
	*User
}

func (c *ApplicationPage) Init() {}

func (c *ApplicationPage) ControllerName() string {
	return "app.controllers.web.application.page"
}

func (c *ApplicationPage) ApplicationPageHandler(rw http.ResponseWriter, r *http.Request) {
}
