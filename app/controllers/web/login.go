package controllers

import (
	"net/http"
)

type Login struct {
	*User
	RedirectAfter string
}

func (c *Login) Init() {
	c.RedirectAfter = "/app"
}

func (c *Login) ControllerName() string {
	return "app.controllers.web.login"
}

func (c *Login) Login(rw http.ResponseWriter, r *http.Request) {
	de := c.Base.JSON.NewDecoder(r.Body)
	var values map[string]string
	de.Decode(&values)

	// Initialize values
	var (
		username string
		pass     string
		exists   bool
	)

	if username, exists = values["Username"]; !exists {
		http.Error(rw, c.Base.Http.Error(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if pass, exists = values["Password"]; !exists {
		http.Error(rw, c.Base.Http.Error(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, err := c.Load(username)

	if err != nil {
		http.Error(rw, c.Base.Http.Error(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if user.Password != pass {
		http.Error(rw, c.Base.Http.Error(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	c.Base.SetSessionKey("client-logged", username, rw, r)
	http.Error(rw, c.Base.Http.Error(http.StatusAccepted), http.StatusAccepted)
}

func (c *Login) Logout(rw http.ResponseWriter, r *http.Request) {
	c.Base.InvalidateSessionKey("client-logged", rw, r)
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}
