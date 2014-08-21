package controllers

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

type HostmaterDB struct {
	Base *a.AppBase
}

func (c *HostmaterDB) CreateInstallation(templ *models.InstallTemplate) {

}

func (c *HostmaterDB) RemoveInstallation() {

}

func (c *HostmaterDB) InstallationExists(name string, rootFolder string) (bool, error) {
	q := "SELECT id FROM installation i WHERE i.name = ? AND i.root_folder = ?"
	row := c.Base.DataStore.DB.QueryRow(q, name, rootFolder)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		log.Println(err)
		return false, err
	}

	if id != 0 {
		return true
	}

	return false
}

func (c *HostmaterDB) CreateSite() {

}

func (c *HostmaterDB) SiteExists() bool {

	return true
}
