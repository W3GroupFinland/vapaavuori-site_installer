package controllers

import (
	"errors"
	//a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

type System struct {
	*Site
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

func (c *System) GetDrupalPlatforms() error {
	pd := c.Base.Config.Platform.Directory
	if pd == "" {
		return errors.New("Platform directory has to be set to get platform listing.")
	}

	files, err := ioutil.ReadDir(pd)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			path := filepath.Join(pd, file.Name())

			exists, info, err := c.InstallRootStatus(path)
			if err != nil {
				return err
			}

			if exists {
				log.Println(info)
			}
		}
	}

	return nil
}
