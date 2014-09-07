package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestDomainsMapToConfigs(t *testing.T) {

	const (
		expectedSite      = "local.bugster2.com"
		expectedSSLCount  = 1
		expectedHTTPCount = 1
		expectedSSLSN     = "local.ssl-bugster2.com"
		expectedSSLHost   = "127.0.0.1"
		expectedHTTPSN    = "local.bugster2.com"
		expectedHTTPHost  = "127.0.0.1"
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
		// Get the configs from database.
		configs, err := a.Application.Controllers.HostMasterDB.GetSiteServerConfigs(site.Id)
		if err != nil {
			t.Error(err)
		}

		// Append them to site config slice.
		site.AddConfigs(configs...)

		// Get site domains from database.
		domains, err := a.Application.Controllers.HostMasterDB.GetSiteDomains(site.Id)
		if err != nil {
			t.Error(err)
		}

		// Append them to site domain slice.
		site.AddDomains(domains...)

		// Map site domains to site configs.
		site.MapDomainsToConfigs(site.Domains)

		sitesByName[site.Name] = site
	}

	if _, exists := sitesByName[expectedSite]; !exists {
		t.Errorf("Site %v should exist in database.", expectedSite)
		return
	}

	if len(sitesByName[expectedSite].SSLConfigs) != expectedSSLCount {
		t.Errorf("Site should have %v SSL server config.", expectedSSLCount)
		return
	}

	if len(sitesByName[expectedSite].HTTPConfigs) != expectedHTTPCount {
		t.Errorf("Site should have %v HTTP server config.", expectedHTTPCount)
		return
	}

	_, err = sites[0].SSLConfigs[0].GetDomainByName(expectedSSLSN, expectedSSLHost)
	if err != nil {
		t.Errorf("Expected: sites[0].SSLConfigs[0].GetDomainByName(%v, %v), Got error: %v",
			expectedSSLSN, expectedSSLHost, err.Error())
	}

	_, err = sites[0].HTTPConfigs[0].GetDomainByName(expectedHTTPSN, expectedHTTPHost)
	if err != nil {
		t.Errorf("Expected: sites[0].HTTPConfigs[0].GetDomainByName(%v, %v), Got error: %v",
			expectedHTTPSN, expectedHTTPHost, err.Error())
	}
}
