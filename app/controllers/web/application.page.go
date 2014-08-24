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
	/*msg := "Hello world!"
	_, err := rw.Write([]byte(msg))
	if err != nil {
		http.Error(rw, "Internal server error.", 500)
	}*/

	// TODO: Remove when finished. Temporary redirecting to html file for developing..
	http.Redirect(rw, r, "/files/app/main.html", http.StatusSeeOther)
}
