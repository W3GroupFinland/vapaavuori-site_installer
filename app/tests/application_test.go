package tests

import (
	"github.com/tuomasvapaavuori/site_installer/app"
	"github.com/tuomasvapaavuori/site_installer/app/tests/hostmaster"
	"github.com/tuomasvapaavuori/site_installer/app/tests/hosts_domains"
	"github.com/tuomasvapaavuori/site_installer/app/tests/system"
	"testing"
)

func TestRunApplicationTests(t *testing.T) {
	// Application initialize.
	app := getApplication()
	// Register application controllers.
	app.RegisterControllers()
	// Close database connection when finished.
	defer app.Base.DataStore.DB.Close()

	// Run tests
	hosts_domains.Init().RunTests(t)
	hostmaster.Init(app).RunTests(t)
	system.Init(app).RunTests(t)
}

func getApplication() *app.Application {
	application := app.Init([]byte(ApplicationConfig1))
	return application
}
