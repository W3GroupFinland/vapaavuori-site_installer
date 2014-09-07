package hostmaster

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (a *ApplicationTests) TestGetSiteDomains(t *testing.T) {

	const (
		expectedSite = "local.bugster2.com"
		expectedDN   = "local.bugster2.com"
		expectedHost = "127.0.0.1"
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
		domains, err := a.Application.Controllers.HostMasterDB.GetSiteDomains(site.Id)
		if err != nil {
			t.Error(err)
		}

		for _, domain := range domains {
			site.AddDomains(domain)
		}

		sitesByName[site.Name] = site
	}

	if _, exists := sitesByName[expectedSite]; !exists {
		t.Errorf("Site %v should exist in database.", expectedSite)
		return
	}

	_, err = sitesByName[expectedSite].GetDomainByName(expectedDN, expectedHost)
	if err != nil {
		t.Errorf("Expected domain name %v on host %v to exist on site %v.\n",
			expectedDN,
			expectedHost,
			expectedSite,
		)
	}
}
