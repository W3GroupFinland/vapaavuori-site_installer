package controllers

import (
	"database/sql"
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

func (c *HostMasterDB) CreateSite(tmpl *models.InstallTemplate, sp *models.SubProcess) (int64, error) {
	sp.Start()

	var id int64
	// Checking if platform exists. If not return error.
	exists, err := c.PlatformIdExists(tmpl.InstallInfo.PlatformId)
	if err != nil {
		return id, err
	}

	if !exists {
		return id, errors.New("Platform doesn't exist.")
	}

	// Check if site exists on given platform.
	exists, err = c.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.SubDirectory)
	if err != nil {
		return id, err
	}

	if exists {
		return id, errors.New("Given site exists already on given platform.")
	}

	// Insert query.
	q := "INSERT INTO site (platform_id,name,db_name,db_user,sub_directory,install_type,template) "
	q += "VALUES(?,?,?,?,?,?,?)"

	res, err := c.Base.DataStore.DB.Exec(q,
		tmpl.InstallInfo.PlatformId,
		tmpl.InstallInfo.SiteName,
		tmpl.DatabaseName.Value,
		tmpl.MysqlUser.Value,
		tmpl.InstallInfo.SubDirectory,
		tmpl.InstallInfo.InstallType,
		tmpl.InstallInfo.TemplatePath)

	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}
	// Set site id to template install info.
	tmpl.InstallInfo.SiteId = id
	// Add rollback functionality.
	tmpl.RollBack.AddDBIdFunction(c.RemoveSiteWithId, id)

	sp.Finish()
	return id, err
}

func (c *HostMasterDB) RemoveSiteWithId(id int64) error {
	q := "DELETE FROM site WHERE id = ?"
	_, err := c.Base.DataStore.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *HostMasterDB) SiteExists(platformId int64, subDirectory string) (bool, error) {
	q := "SELECT id FROM site s WHERE s.platform_id = ? AND s.sub_directory = ?"
	row := c.Base.DataStore.DB.QueryRow(q, platformId, subDirectory)

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
