package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func getApplication() *app.Application {
	application := app.Init([]byte(ApplicationConfig1))
	return application
}

func RunTests(t *testing.T) {
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
	// Test site creation, database server configs and domains.
	appTests.TestCreateConfigsAndDomains(t)
}

type ApplicationTests struct {
	Application *app.Application
}

func (a *ApplicationTests) RandomizeDatabaseValues(tmpl *models.InstallTemplate) {
	tmpl.MysqlUser = models.RandomValue{Random: true}
	tmpl.MysqlPassword = models.RandomValue{Random: true}
	tmpl.MysqlUserHosts = models.MysqlUserHosts{Hosts: []string{"127.0.0.1", "localhost"}}
	tmpl.DatabaseName = models.RandomValue{Random: true}
	tmpl.MysqlUserPrivileges = models.MysqlUserPrivileges{Privileges: []string{"ALL"}}
	tmpl.MysqlGrantOption = models.MysqlGrantOption{Value: true}
}
