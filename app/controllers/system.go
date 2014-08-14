package controllers

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"log"
	"os/exec"
)

type System struct {
	Base *a.AppBase
}

func (c *System) ApacheRestart() error {
	out, err := exec.Command("apachectl", "restart").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(string(out))

	return nil
}
