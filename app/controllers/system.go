package controllers

import (
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

type System struct {
	Site       *Site
	HostMaster *HostMasterDB
}

func (c *System) HttpServerRestart() error {
	if c.Site.Base.Commands.HttpServer.Restart.Command == "" {
		return errors.New("No command set.")
	}
	cmd := c.Site.Base.Commands.HttpServer.Restart
	out, err := exec.Command(cmd.Command, cmd.Arguments...).Output()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(string(out))

	return nil
}

func (c *System) GetDrupalPlatforms() (models.PlatformList, error) {
	var platforms models.PlatformList

	pd := c.Site.Base.Config.Platform.Directory
	if pd == "" {
		return platforms, errors.New("Platform directory has to be set to get platform listing.")
	}

	files, err := ioutil.ReadDir(pd)
	if err != nil {
		return platforms, err
	}

	for _, file := range files {
		if file.IsDir() {
			path := filepath.Join(pd, file.Name())

			exists, info, err := c.Site.InstallRootStatus(path)
			if err != nil {
				return platforms, err
			}

			if !exists {
				continue
			}

			platform := models.PlatformInfo{
				RootInfo: info,
				Name:     file.Name(),
			}

			// Check if platform is already registered.
			exists, id, err := c.HostMaster.PlatformExists(file.Name(), pd)
			if err != nil {
				return platforms, err
			}

			if exists {
				platform.Registered = true
				platform.PlatformId = id
			}

			// Value is not given out for security reasons.
			// Config platform directory is used instead when creating platforms.
			platform.RootInfo.DrupalRoot = "-- Obfuscated --"

			platforms.Add(pd, &platform)
		}
	}

	return platforms, nil
}
