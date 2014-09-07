package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestGetPlatformSites(t *testing.T) {

	const (
		expectedSite = "local.bugster2.com"
	)

	tmpl := a.TestCreateConfigsAndDomains_Data(t)
	sp := a.GetTestSubProcessChannel()

	defer tmpl.RollBack.Execute()
	a.CreateConfigAndDomains(tmpl, sp, t)

	sites, err := a.Application.Controllers.HostMasterDB.GetPlatformSites(tmpl.InstallInfo.PlatformId)
	if err != nil {
		t.Error(err)
	}

	var sitesByName = make(map[string]*models.SiteInfo)
	for _, site := range sites {
		sitesByName[site.Name] = site
	}

	if _, exists := sitesByName[expectedSite]; !exists {
		t.Errorf("Site %v should exist in database.", expectedSite)
	}
}
