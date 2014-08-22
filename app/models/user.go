package models

import (
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
)

type RandomValue struct {
	Value  string
	Random bool
}

func (rv *RandomValue) Randomize(length int) {
	rv.Value = utils.RandomString(length)
}

type User struct {
	Uid      int64
	Username string
	Password string
	Mail     string
	Status   bool
}

type UserSend struct {
	Username string
	Mail     string
	Status   bool
}

const (
	USERNAME_LENGHT = 20
	PASSWORD_LENGTH = 20
)

var (
	NoUserFoundError = errors.New("No user found.")
)

func (u *User) ValueReferences() (*int64, *string, *string, *string, *bool) {
	return &u.Uid, &u.Username, &u.Mail, &u.Password, &u.Status
}
