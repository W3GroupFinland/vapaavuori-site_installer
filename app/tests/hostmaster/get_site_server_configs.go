package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestGetSiteServerConfigs(t *testing.T) {

	const (
		expectedSite      = "local.bugster2.com"
		expectedSSLCount  = 1
		expectedHTTPCount = 1
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
		configs, err := a.Application.Controllers.HostMasterDB.GetSiteServerConfigs(site.Id)
		if err != nil {
			t.Error(err)
		}

		for _, config := range configs {
			site.AddConfigs(config)
		}

		sitesByName[site.Name] = site
	}

	if _, exists := sitesByName[expectedSite]; !exists {
		t.Errorf("Site %v should exist in database.", expectedSite)
		return
	}

	if len(sitesByName[expectedSite].SSLConfigs) != expectedSSLCount {
		t.Errorf("Site should have %v SSL server config.", expectedSSLCount)
	}

	if len(sitesByName[expectedSite].HTTPConfigs) != expectedHTTPCount {
		t.Errorf("Site should have %v HTTP server config.", expectedHTTPCount)
	}
}
