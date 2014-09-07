package controllers

import (
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
)

func (c *HostMasterDB) CreateServerConfigs(tmpl *models.InstallTemplate, sp *models.SubProcess) error {
	sp.Start()

	if tmpl.HttpServer.Include {
		si := tmpl.HttpServer
		id, err := c.CreateServerConfig(tmpl,
			&models.DatabaseServerConfig{
				SiteId:     tmpl.InstallInfo.SiteId,
				ServerType: si.Type,
				Template:   si.Template,
				Port:       si.Port,
				ConfigRoot: si.ConfigRoot,
				ConfigFile: si.ConfigFile,
				Secured:    false,
			})

		if err != nil {
			return err
		}

		tmpl.HttpServer.ServerConfigId = id

		// Helper function to update domain objects.
		// Without site id and server id domain can't be created to database.
		err = tmpl.HttpServer.UpdateDomainIds(tmpl.InstallInfo.SiteId)
		if err != nil {
			return err
		}
	}

	if tmpl.SSLServer.Include {
		si := tmpl.SSLServer
		id, err := c.CreateServerConfig(tmpl,
			&models.DatabaseServerConfig{
				SiteId:     tmpl.InstallInfo.SiteId,
				ServerType: si.Type,
				Template:   si.Template,
				Port:       si.Port,
				ConfigRoot: si.ConfigRoot,
				ConfigFile: si.ConfigFile,
				Secured:    true,
				Cert:       si.Certificate,
				CertKey:    si.Key,
			})

		if err != nil {
			return err
		}

		tmpl.SSLServer.ServerConfigId = id

		// Helper function to update domain objects.
		// Without site id and server id domain can't be created to database.
		err = tmpl.SSLServer.UpdateDomainIds(tmpl.InstallInfo.SiteId)
		if err != nil {
			return err
		}
	}

	sp.Finish()
	return nil
}

func (c *HostMasterDB) CreateServerConfig(tmpl *models.InstallTemplate, cs *models.DatabaseServerConfig) (int64, error) {
	if cs.SiteId == 0 {
		return 0, errors.New("No site id when creating server config.")
	}

	var id int64

	q := "INSERT INTO server_config (site_id, server_type, template, port, config_root, config_file, secured, cert, cert_key) "
	q += "VALUES(?,?,?,?,?,?,?,?,?)"

	res, err := c.Base.DataStore.DB.Exec(q,
		cs.SiteId,
		cs.ServerType,
		cs.Template,
		cs.Port,
		cs.ConfigRoot,
		cs.ConfigFile,
		cs.Secured,
		cs.Cert,
		cs.CertKey)

	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	// Add rollback functionality.
	tmpl.RollBack.AddDBIdFunction(c.RemoveServerConfigWithId, id)

	return id, nil
}

func (c *HostMasterDB) RemoveServerConfigWithId(id int64) error {
	q := "DELETE FROM server_config WHERE id = ?"
	_, err := c.Base.DataStore.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *HostMasterDB) GetSiteServerConfigs(siteId int64) ([]*models.DatabaseServerConfig, error) {
	var configs []*models.DatabaseServerConfig
	q := "SELECT sc.* FROM site s "
	q += "JOIN server_config sc ON sc.site_id = s.id "
	q += "WHERE s.id = ?"

	rows, err := c.Base.DataStore.DB.Query(q, siteId)
	if err != nil {
		return configs, err
	}

	defer rows.Close()

	for rows.Next() {
		c := models.DatabaseServerConfig{}
		err := rows.Scan(
			&c.Id,
			&c.SiteId,
			&c.ServerType,
			&c.Template,
			&c.Port,
			&c.ConfigRoot,
			&c.ConfigFile,
			&c.Secured,
			&c.Cert,
			&c.CertKey,
		)
		if err != nil {
			return configs, err
		}

		configs = append(configs, &c)
	}

	if err := rows.Err(); err != nil {
		return configs, err
	}

	return configs, nil
}
