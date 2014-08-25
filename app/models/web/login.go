package web_models

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
)

type Login struct {
	Username string
	Password string
	Messages *FormMessages
}

type FormMessages struct {
	Messages map[string]string
	Errors   map[string]string
}

func (l *Login) Validate(user *models.User, err error) bool {
	l.Messages.Errors = make(map[string]string)

	if err != nil || user.Password != l.Password {
		l.Messages.Errors["General"] = "Please check username and password!"

		return false
	}

	return true
}
