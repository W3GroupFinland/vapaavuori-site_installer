package hostmaster

import (
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestCreateConfigsAndDomains(t *testing.T) {
	tmpl := a.TestCreateConfigsAndDomains_Data(t)
	defer tmpl.RollBack.Execute()

	// Create new platform
	t.Log("Creating new platform to database..")
	_, err := a.Application.Controllers.HostMasterDB.CreatePlatform(tmpl)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Creating new database..")
	di, err := a.Application.Controllers.Site.CreateDatabase(tmpl)
	if err != nil {
		t.Error(err)
		return
	}

	// Check database exists.
	t.Log("Checking database exists..")
	if !a.Application.Base.DataStore.CheckDatabaseExists(di.DbName.Value) {
		t.Errorf("Failed creating database %v.\n", di.DbName.Value)
		return
	}

	// Add site to platform.
	t.Log("Adding new site to platform..")
	_, err = a.Application.Controllers.HostMasterDB.CreateSite(tmpl)
	if err != nil {
		t.Error(err)
		return
	}

	// Check that site id is not zero.
	t.Log("Checking that site id is not zero..")
	if tmpl.InstallInfo.SiteId == 0 {
		t.Error("Site id should not be zero after creating site.")
		return
	}

	// Check that new site exists in database.
	t.Log("Check that new site exists in hostmaster database..")
	exists, err := a.Application.Controllers.HostMasterDB.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.SubDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	if !exists {
		t.Errorf("Site should exist in database.")
		return
	}

	t.Log("Creating server configs to hostmaster database..")
	err = a.Application.Controllers.HostMasterDB.CreateServerConfigs(tmpl)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Checking server domains has ids.")
	err = a.CheckServerDomainsHasIds(tmpl, t)
	if err != nil {
		t.Error(err)
		return
	}

	// Get domains from template.
	t.Log("Get domains from install template.")
	domains := a.Application.Controllers.Site.GetSiteTemplateDomains(tmpl)

	// Create site domains.
	t.Log("Creating site domains.")
	err = a.Application.Controllers.HostMasterDB.CreateSiteDomains(tmpl, domains)
	if err != nil {
		t.Error(err)
	}
}

func (a *ApplicationTests) CheckServerDomainsHasIds(tmpl *models.InstallTemplate, t *testing.T) error {
	if tmpl.HttpServer.Template != "" {
		err := a.CheckDomainHasIds(tmpl.HttpServer.DomainInfo, t)
		if err != nil {
			return err
		}

		for _, alias := range tmpl.HttpServer.DomainAliases {
			err := a.CheckDomainHasIds(alias, t)
			if err != nil {
				return err
			}
		}
	}
	if tmpl.SSLServer.Template != "" {
		err := a.CheckDomainHasIds(tmpl.SSLServer.DomainInfo, t)
		if err != nil {
			return err
		}

		for _, alias := range tmpl.SSLServer.DomainAliases {
			err := a.CheckDomainHasIds(alias, t)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *ApplicationTests) CheckDomainHasIds(domain *models.Domain, t *testing.T) error {
	var msg string
	if domain.SiteId == 0 {
		msg += "Site id is zero. "
	}
	if domain.ServerConfigId == 0 {
		msg += "Server config id is zero."
	}

	if msg != "" {
		return errors.New(msg)
	}

	t.Logf("Domain %v has proper site id (%v) and server config id (%v)",
		domain.DomainName, domain.SiteId, domain.ServerConfigId)

	return nil
}

func (a *ApplicationTests) TestCreateConfigsAndDomains_Data(t *testing.T) *models.InstallTemplate {
	if a.Application.Base.Config.Backup.Directory == "" {
		t.Error("Backup directory must be set for test to work.")
		return &models.InstallTemplate{}
	}

	tmpl := models.InstallTemplate{
		InstallInfo: models.SiteInstallInfo{
			PlatformName:     "bugsters",
			DrupalRoot:       "/bugsters/gatan",
			InstallType:      "test-install",
			TemplatePath:     "/bug/it",
			SiteName:         "local.bugster2.com",
			SubDirectory:     "local.bugster.com",
			HttpUser:         "_www",
			HttpGroup:        "_www",
			ServerConfigRoot: a.Application.Base.Config.Backup.Directory,
		},
		HttpServer: models.HttpServerTemplate{
			Type:       "apache",
			Template:   "/bug/it/default.conf",
			Port:       8888,
			ConfigRoot: "/bug/it/apache/conf",
			ConfigFile: "http.local.bugster.com",
		},
		SSLServer: models.SSLServerTemplate{
			Certificate: "/bug/it/certs/bugster.crt",
			Key:         "/bug/it/certs/bugster.key",
		},
	}

	tmpl.SSLServer.Type = "apache"
	tmpl.SSLServer.Template = "/bug/it/ssl-default.conf"
	tmpl.SSLServer.Port = 443
	tmpl.SSLServer.ConfigRoot = "/bug/it/apache/conf"
	tmpl.SSLServer.ConfigFile = "ssl.local.bugster.com"

	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)

	// Regular server domains.
	tmpl.HttpServer.DomainInfo = &models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.bugster2.com",
		Type:       models.DomainTypeServerName,
	}
	tmpl.HttpServer.DomainAliases = []*models.Domain{
		&models.Domain{
			Host:       "127.0.0.1",
			DomainName: "local.alias.bugster2.com",
			Type:       models.DomainTypeServerAlias,
		},
		&models.Domain{
			Host:       "127.0.0.1",
			DomainName: "local.alias2.bugster2.com",
			Type:       models.DomainTypeServerAlias,
		},
	}
	// SSL Domains
	tmpl.SSLServer.DomainInfo = &models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.ssl-bugster2.com",
		Type:       models.DomainTypeServerName,
	}
	tmpl.SSLServer.DomainAliases = []*models.Domain{
		&models.Domain{
			Host:       "127.0.0.1",
			DomainName: "local.ssl-alias-bugster2.com",
			Type:       models.DomainTypeServerAlias,
		},
		&models.Domain{
			Host:       "127.0.0.1",
			DomainName: "local.ssl-alias2-bugster2.com",
			Type:       models.DomainTypeServerAlias,
		},
	}

	a.RandomizeDatabaseValues(&tmpl)

	return &tmpl
}
