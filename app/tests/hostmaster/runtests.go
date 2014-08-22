package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

type ApplicationTests struct {
	Application *app.Application
}

func Init(a *app.Application) *ApplicationTests {
	return &ApplicationTests{Application: a}
}

func (a *ApplicationTests) RunTests(t *testing.T) {
	// Test reads config.
	a.TestReadsConfig(t)
	// Test creates database.
	a.TestCreatesDatabase(t)
	// Test platform exists.
	a.TestPlatformExists(t)
	// Test site creation.
	a.TestCreateSite(t)
	// Test site creation, database server configs and domains.
	a.TestCreateConfigsAndDomains(t)
}

func (a *ApplicationTests) RandomizeDatabaseValues(tmpl *models.InstallTemplate) {
	tmpl.MysqlUser = models.RandomValue{Random: true}
	tmpl.MysqlPassword = models.RandomValue{Random: true}
	tmpl.MysqlUserHosts = models.MysqlUserHosts{Hosts: []string{"127.0.0.1", "localhost"}}
	tmpl.DatabaseName = models.RandomValue{Random: true}
	tmpl.MysqlUserPrivileges = models.MysqlUserPrivileges{Privileges: []string{"ALL"}}
	tmpl.MysqlGrantOption = models.MysqlGrantOption{Value: true}
}
