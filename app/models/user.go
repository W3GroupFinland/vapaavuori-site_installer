package models

import (
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
)

type User struct {
	Username string
	Password string
}

type RandomValue struct {
	Value  string
	Random bool
}

func (rv *RandomValue) Randomize(length int) {
	rv.Value = utils.RandomString(length)
}
