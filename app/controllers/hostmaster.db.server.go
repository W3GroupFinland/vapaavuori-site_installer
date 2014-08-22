package controllers

import (
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
)

func (c *HostMasterDB) CreateServerConfigs(tmpl *models.InstallTemplate) {
	if tmpl.HttpServer.Template != "" {
		si := tmpl.HttpServer
		id, err := c.CreateServerConfig(tmpl,
			&models.DatabaseServerConfig{
				SiteId:     tmpl.InstallInfo.SiteId,
				ServerType: si.Type,
				Template:   si.Template,
				Port:       si.Port,
				ConfigRoot: si.ConfigRoot,
				ConfigFile: si.Type,
				Secured:    false,
			})

		if err != nil {
			return err
		}

		templ.HttpServer.ServerId = id

		// Helper function to update domain objects.
		// Without site id and server id domain can't be created to database.
		c.UpdateServerDomains(
			tmpl.InstallInfo.SiteId,
			templ.HttpServer.ServerId,
			tmpl.HttpServer.DomainInfo,
			tmpl.HttpServer.DomainAliases)
	}

	if tmpl.SSLServer.Template != "" {
		si := tmpl.SSLServer
		id, err := c.CreateServerConfig(tmpl,
			&models.DatabaseServerConfig{
				SiteId:     tmpl.InstallInfo.SiteId,
				ServerType: si.Type,
				Template:   si.Template,
				Port:       si.Port,
				ConfigRoot: si.ConfigRoot,
				ConfigFile: si.Type,
				Secured:    true,
				Cert:       si.Certificate,
				CertKey:    si.Key,
			})

		if err != nil {
			return err
		}

		templ.SSLServer.ServerId = id

		// Helper function to update domain objects.
		// Without site id and server id domain can't be created to database.
		c.UpdateServerDomains(
			tmpl.InstallInfo.SiteId,
			templ.SSLServer.ServerId,
			tmpl.SSLServer.DomainInfo,
			tmpl.SSLServer.DomainAliases)
	}
}

// Function to update site id and server config id to domains.
func (c *HostMasterDB) UpdateServerDomains(siteId int64,
	serverConfigId int64,
	serverName *models.Domain,
	serverAliases []*models.Domain) {

	serverName.SiteId = siteId
	serverName.ServerConfigId = serverConfigId

	for _, domain := range serverAliases {
		domain.SiteId = siteId
		domain.ServerConfigId = serverConfigId
	}
}

func (c *HostMasterDB) CreateServerConfig(tmpl *models.InstallTemplate, cs *models.DatabaseServerConfig) (int64, error) {
	if config.SiteId == 0 {
		return 0, errors.New("No site id when creating server config.")
	}

	q := "INSERT INTO site (site_id, server_type, template, port, config_root, config_file, secured, cert, cert_key) "
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
