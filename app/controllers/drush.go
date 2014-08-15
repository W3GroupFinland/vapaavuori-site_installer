package controllers

import (
	"bytes"
	"errors"
	"fmt"
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"os/exec"
)

const (
	NoCommandFound = "No drush command found."
	DrushCommand   = "drush"

	// Drush arguments
	DrushDatabaseUrlArg = "--db-url"
	DrushSiteNameArg    = "--site-name"
	DrushSiteSubDirArg  = "--sites-subdir"
)

type Drush struct {
	Base      *a.AppBase
	DrushInfo *models.DrushInfo
}

func (d *Drush) Init() {
	d.DrushInfo = &models.DrushInfo{}

	path, err := d.Which()
	if err != nil {
		msg := NoCommandFound + " You have to install drush before continue.\n"
		msg += "Error: %v\n"
		msg += "Path: %v\n"
		log.Fatalf(msg, err.Error(), path)
	}

	d.DrushInfo.Executable = path
	version, err := d.ReadVersion()
	if err != nil {
		log.Fatalln(err)
	}

	d.DrushInfo.Version = version
}

func (d *Drush) ReadVersion() (string, error) {
	out, err := exec.Command(d.DrushInfo.Executable, "--version").Output()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(out), nil
}

func (d *Drush) Which() (string, error) {
	path := d.Base.Config.Drush.Path

	// If drush path was not given in config-file, try to search
	// command with which.
	if path == "" {
		out, err := exec.Command("which", DrushCommand).Output()
		if err != nil {
			log.Println(err)
		}

		bytes := utils.StripPathWhiteSpace(out)
		if len(bytes) == 0 {
			return "", errors.New(NoCommandFound)
		}

		path = string(bytes)
	}

	if !utils.FileExists(path) {
		return "", errors.New(NoCommandFound)
	}

	return path, nil
}

func (d *Drush) Run(args ...string) (string, error) {
	cmd := exec.Command("drush", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	// Run command.
	err := cmd.Run()
	// Check for errors.
	if err != nil {
		log.Println(err)
		return out.String(), err
	}

	log.Println(out.String())

	return out.String(), nil
}

func (d *Drush) ArgumentStringFormat(name string) string {
	var str string
	switch name {
	case DrushDatabaseUrlArg:
		str = "--db-url=mysql://%v:%v@%v:%v/%v"
		break
	case DrushSiteNameArg:
		str = "--site-name=%v"
		break
	case DrushSiteSubDirArg:
		str = "--sites-subdir=%v"
		break
	}

	return str
}

func (d *Drush) FormatDatabaseStr(db *database.DatabaseInfo) string {
	return fmt.Sprintf(d.ArgumentStringFormat(DrushDatabaseUrlArg),
		db.User.Value,
		db.Password.Value,
		d.Base.Config.Mysql.Host,
		d.Base.Config.Mysql.Port,
		db.DbName.Value)
}

func (d *Drush) FormatSiteNameStr(ic *models.SiteInstallConfig) string {
	return fmt.Sprintf(d.ArgumentStringFormat(DrushSiteNameArg), ic.SiteName)
}

func (d *Drush) FormatSiteSubDirStr(ic *models.SiteInstallConfig) string {
	return fmt.Sprintf(d.ArgumentStringFormat(DrushSiteSubDirArg), ic.SubDirectory)
}
