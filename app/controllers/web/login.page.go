package controllers

import (
	"github.com/tuomasvapaavuori/site_installer/app/models/web"
	"log"
	"net/http"
)

func (c *Login) Page(rw http.ResponseWriter, r *http.Request) {
	_, valid := c.Current(rw, r)
	// If user already logged in redirect to application page.
	if valid {
		http.Redirect(rw, r, c.RedirectAfter, http.StatusSeeOther)
		return
	}

	err := c.Base.Templates.Templates.ExecuteTemplate(rw, "login.html", nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Posts login form data and logs user in.
func (c *Login) LoginPostHandler(rw http.ResponseWriter, r *http.Request) {
	// Create struct of form data to handle possible errors.
	formData := web_models.Login{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Messages: &web_models.FormMessages{}}

	user, err := c.Load(formData.Username)
	log.Println(user)
	if err != nil {
		log.Println(err)
	}

	if formData.Validate(user, err) == false {
		// If errors found render form with errors.
		err := c.Base.Templates.Templates.ExecuteTemplate(rw, "login.html", formData)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Set user session.
		c.Base.SetSessionKey("client-logged", user.Username, rw, r)
		http.Redirect(rw, r, c.RedirectAfter, http.StatusSeeOther)
	}
}
