package tests

import (
	"database/sql"
	"github.com/tuomasvapaavuori/site_installer/app"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func getApplication() *app.Application {
	application := app.Init([]byte(ApplicationConfig1))
	return application
}

func TestRunApplicationTests(t *testing.T) {
	// Application initialize.
	app := getApplication()
	app.RegisterControllers()
	defer app.Base.DataStore.DB.Close()

	appTests := ApplicationTests{app}
	// Test reads config.
	appTests.TestReadsConfig(t)
	// Test creates database.
	appTests.TestCreatesDatabase(t)
	// Test platform exists.
	appTests.TestPlatformExists(t)
	// Test site creation.
	appTests.TestCreateSite(t)
}

type ApplicationTests struct {
	Application *app.Application
}

func (a *ApplicationTests) TestReadsConfig(t *testing.T) {
	// If application mysql user is empty test fails.
	if a.Application.Base.Config.Mysql.User == "" {
		t.Error("app.Base.Config.Mysql.User == \"\", expected not empty value.")
	}
}

func (a *ApplicationTests) RandomizeDatabaseValues(tmpl *models.InstallTemplate) {
	tmpl.MysqlUser = models.RandomValue{Random: true}
	tmpl.MysqlPassword = models.RandomValue{Random: true}
	tmpl.MysqlUserHosts = models.MysqlUserHosts{Hosts: []string{"127.0.0.1", "localhost"}}
	tmpl.DatabaseName = models.RandomValue{Random: true}
	tmpl.MysqlUserPrivileges = models.MysqlUserPrivileges{Privileges: []string{"ALL"}}
	tmpl.MysqlGrantOption = models.MysqlGrantOption{Value: true}
}

func (a *ApplicationTests) TestCreatesDatabase(t *testing.T) {
	tmpl := models.InstallTemplate{}
	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)
	a.RandomizeDatabaseValues(&tmpl)

	di, err := a.Application.Controllers.Site.CreateDatabase(&tmpl)
	if err != nil {
		t.Error(err)
	}

	if !a.Application.Base.DataStore.CheckDatabaseExists(di.DbName.Value) {
		t.Errorf("Failed creating database %v.\n", di.DbName.Value)
	}

	tmpl.RollBack.Execute()

	if a.Application.Base.DataStore.CheckDatabaseExists(di.DbName.Value) {
		t.Errorf("Database rollback failed on database %v.\n", di.DbName.Value)
	}
}

func (a *ApplicationTests) TestPlatformExists(t *testing.T) {
	const (
		installRoot = "testing"
		installName = "testing"
	)
	// Database should be empty so no rows exists.
	exists, err := a.Application.Controllers.HostMasterDB.PlatformExists(installName, installRoot)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("Platform with name %v and install root %v should not exist.", installName, installRoot)
	}

	tmpl := models.InstallTemplate{InstallInfo: models.SiteInstallInfo{PlatformName: "bugster", DrupalRoot: "/bugsters/gatan"}}
	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)

	_, err = a.Application.Controllers.HostMasterDB.CreatePlatform(&tmpl)
	if err != nil {
		t.Error(err)
	}

	// Database should be empty so no rows exists.
	exists, err = a.Application.Controllers.HostMasterDB.PlatformExists(tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("Platform with name %v and install root %v should exist.", tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	}

	// Removing inserted lines.
	err = a.Application.Controllers.HostMasterDB.RemovePlatform(tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	if err != nil {
		t.Error(err)
	}

	// No install should exist.
	exists, err = a.Application.Controllers.HostMasterDB.PlatformExists(tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("Platform with name %v and install root %v should not exist.", tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	}
}

func (a *ApplicationTests) TestCreateSite(t *testing.T) {
	tmpl := models.InstallTemplate{
		InstallInfo: models.SiteInstallInfo{PlatformName: "bugster", DrupalRoot: "/bugsters/gatan"},
	}

	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)

	// Create new platform
	_, err := a.Application.Controllers.HostMasterDB.CreatePlatform(&tmpl)
	if err != nil {
		t.Error(err)
	}

	// Check does site exist in database.
	exists, err := a.Application.Controllers.HostMasterDB.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.PlatformName)
	if err != nil {
		t.Error(err)
	}

	if exists {
		t.Errorf("Site shouldn't exist in database.")
		return
	}

	a.RandomizeDatabaseValues(&tmpl)

	di, err := a.Application.Controllers.Site.CreateDatabase(&tmpl)
	if err != nil {
		t.Error(err)
	}

	// Check database exists.
	if !a.Application.Base.DataStore.CheckDatabaseExists(di.DbName.Value) {
		t.Errorf("Failed creating database %v.\n", di.DbName.Value)
	}

	// Add site to platform.
	_, err = a.Application.Controllers.HostMasterDB.CreateSite(&tmpl)
	if err != nil {
		t.Error(err)
	}

	// Check that new site exists in database.
	exists, err = a.Application.Controllers.HostMasterDB.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.SubDirectory)
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Errorf("Site should exist in database.")
		return
	}

	tmpl.RollBack.Execute()
}

func (a *ApplicationTests) TestCreateDomains(tmpl *models.InstallTemplate, t *testing.T) {
	domains := models.SiteDomains{
		SiteName:     tmpl.InstallInfo.SiteName,
		SubDirectory: tmpl.InstallInfo.SubDirectory,
		Domains:      make(map[string]*models.Domain),
	}

	domains.SetDomain(&models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.bugster.com",
		Type:       models.DomainTypeServerName,
	})

	domains.SetDomain(&models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.ssl-bugster.com",
		Type:       models.DomainTypeServerName,
	})
}
