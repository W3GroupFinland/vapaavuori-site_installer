package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestCreateSite(t *testing.T) {
	tmpl := models.InstallTemplate{
		InstallInfo: models.SiteInstallInfo{PlatformName: "bugsteeere", DrupalRoot: "/bugsters/gatan"},
	}

	sp := a.GetTestSubProcessChannel()

	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)
	defer tmpl.RollBack.Execute()

	// Create new platform
	_, err := a.Application.Controllers.HostMasterDB.CreatePlatform(&tmpl)
	if err != nil {
		t.Error(err)
		return
	}

	// Check does site exist in database.
	exists, err := a.Application.Controllers.HostMasterDB.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.PlatformName)
	if err != nil {
		t.Error(err)
		return
	}

	if exists {
		t.Errorf("Site shouldn't exist in database.")
		return
	}

	a.RandomizeDatabaseValues(&tmpl)

	di, err := a.Application.Controllers.Site.CreateDatabase(&tmpl)
	if err != nil {
		t.Error(err)
		return
	}

	// Check database exists.
	if !a.Application.Base.DataStore.CheckDatabaseExists(di.DbName.Value) {
		t.Errorf("Failed creating database %v.\n", di.DbName.Value)
		return
	}

	// Add site to platform.
	_, err = a.Application.Controllers.HostMasterDB.CreateSite(&tmpl, sp)
	if err != nil {
		t.Error(err)
		return
	}

	// Check that new site exists in database.
	exists, err = a.Application.Controllers.HostMasterDB.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.SubDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	if !exists {
		t.Errorf("Site should exist in database.")
		return
	}

	// Inform channel that test process is finished.
	// This makes loop exit.
	sp.StateChannel <- models.SubProcessState{State: ProcessStateTestFinished}
}
