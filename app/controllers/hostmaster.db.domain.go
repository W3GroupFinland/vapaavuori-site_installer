package controllers

import (
	"database/sql"
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

func (c *HostMasterDB) CreateSiteDomains(tmpl *models.InstallTemplate, domains *models.SiteDomains) error {
	if tmpl.InstallInfo.SiteId == 0 {
		return errors.New("Site id can't be zero.")
	}

	for _, domain := range domains.Domains {
		exists, err := c.DomainExists(domain)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		_, err = c.CreateSiteDomain(tmpl, domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *HostMasterDB) CreateSiteDomain(tmpl *models.InstallTemplate, domain *models.Domain) (int64, error) {
	var id int64
	if domain.SiteId == 0 || domain.ServerConfigId == 0 {
		return id, errors.New("Site id or server config id is zero.")
	}

	q := "INSERT INTO domain (site_id, server_config_id, type, name, host) VALUES(?,?,?,?,?)"
	res, err := c.Base.DataStore.DB.Exec(q,
		domain.SiteId,
		domain.ServerConfigId,
		domain.Type,
		domain.DomainName,
		domain.Host)

	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	tmpl.RollBack.AddDBIdFunction(c.RemoveSiteDomainWithId, id)

	return id, nil
}

func (c *HostMasterDB) RemoveSiteDomainWithId(id int64) error {
	q := "DELETE FROM domain WHERE id = ?"
	_, err := c.Base.DataStore.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *HostMasterDB) DomainExists(domain *models.Domain) (bool, error) {
	q := "SELECT id FROM domain d WHERE d.name = ? AND d.host = ?"
	row := c.Base.DataStore.DB.QueryRow(q, domain.DomainName, domain.Host)

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
