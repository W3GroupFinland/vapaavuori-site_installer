package controllers

import (
	"errors"
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"log"
	"os/exec"
)

type System struct {
	Base *a.AppBase
}

func (c *System) HttpServerRestart() error {
	if c.Base.Commands.HttpServer.Restart.Command == "" {
		return errors.New("No command set.")
	}
	cmd := c.Base.Commands.HttpServer.Restart
	out, err := exec.Command(cmd.Command, cmd.Arguments...).Output()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(string(out))

	return nil
}
