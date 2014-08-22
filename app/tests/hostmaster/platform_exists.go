package hostmaster

import (
	"database/sql"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestPlatformExists(t *testing.T) {
	const (
		installRoot = "testing"
		installName = "testing"
	)
	// Database should be empty so no rows exists.
	exists, _, err := a.Application.Controllers.HostMasterDB.PlatformExists(installName, installRoot)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("Platform with name %v and install root %v should not exist.", installName, installRoot)
	}

	tmpl := models.InstallTemplate{InstallInfo: models.SiteInstallInfo{PlatformName: "bugsters", DrupalRoot: "/bugsters/gatan"}}
	// Initialize rollback functionality.
	tmpl.RollBack = models.NewSiteRollBack(&tmpl)

	_, err = a.Application.Controllers.HostMasterDB.CreatePlatform(&tmpl)
	if err != nil {
		t.Error(err)
	}

	// Database should be empty so no rows exists.
	exists, _, err = a.Application.Controllers.HostMasterDB.PlatformExists(tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
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
	exists, _, err = a.Application.Controllers.HostMasterDB.PlatformExists(tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if exists {
		t.Errorf("Platform with name %v and install root %v should not exist.", tmpl.InstallInfo.PlatformName, tmpl.InstallInfo.DrupalRoot)
	}
}
