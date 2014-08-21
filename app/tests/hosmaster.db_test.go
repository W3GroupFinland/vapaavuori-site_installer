package tests

import (
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

func (a *ApplicationTests) TestCreatesDatabase(t *testing.T) {
	tmpl := models.InstallTemplate{
		MysqlUser:           models.RandomValue{Random: true},
		MysqlPassword:       models.RandomValue{Random: true},
		MysqlUserHosts:      models.MysqlUserHosts{Hosts: []string{"127.0.0.1", "localhost"}},
		DatabaseName:        models.RandomValue{Random: true},
		MysqlUserPrivileges: models.MysqlUserPrivileges{Privileges: []string{"ALL"}},
		MysqlGrantOption:    models.MysqlGrantOption{Value: true},
	}

	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)

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
