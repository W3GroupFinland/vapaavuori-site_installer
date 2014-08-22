package controllers

import (
	"database/sql"
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

func (c *HostMasterDB) CreatePlatform(tmpl *models.InstallTemplate) (int64, error) {
	var id int64

	exists, err := c.PlatformExists(tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	if err != nil {
		return id, err
	}

	if exists {
		return id, errors.New("Site platform exists already.")
	}

	q := "INSERT INTO platform (name, root_folder) VALUES(?, ?)"
	res, err := c.Base.DataStore.DB.Exec(q, tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	// Add database id to install info.
	tmpl.InstallInfo.PlatformId = id
	// Add rollback functionality.
	tmpl.RollBack.AddDBIdFunction(c.RemovePlatformWithId, id)

	return id, nil
}

func (c *HostMasterDB) RemovePlatform(sitename string, rootfolder string) error {
	q := "DELETE FROM platform WHERE name = ? AND root_folder = ?"
	_, err := c.Base.DataStore.DB.Exec(q, sitename, rootfolder)
	if err != nil {
		return err
	}

	return nil
}

func (c *HostMasterDB) RemovePlatformWithId(id int64) error {
	q := "DELETE FROM platform WHERE platform.id = ?"
	_, err := c.Base.DataStore.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *HostMasterDB) PlatformIdExists(id int64) (bool, error) {
	q := "SELECT name FROM platform p WHERE p.id = ?"
	row := c.Base.DataStore.DB.QueryRow(q, id)

	var name string
	err := row.Scan(&name)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (c *HostMasterDB) PlatformExists(name string, rootFolder string) (bool, error) {
	q := "SELECT id FROM platform p WHERE p.name = ? AND p.root_folder = ?"
	row := c.Base.DataStore.DB.QueryRow(q, name, rootFolder)

	var id int64
	err := row.Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false, err
	}

	if id != 0 {
		return true, nil
	}

	return false, nil
}
