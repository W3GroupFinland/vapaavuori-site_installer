package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

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
