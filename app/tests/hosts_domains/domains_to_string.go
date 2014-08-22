package hosts_domains

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (h *HostsDomains) TestHostsDomainsToString(t *testing.T) {
	const (
		expectedHostStr1  = "127.0.0.1 local.luolamies.fi local.amputaatio.fi local.amputaatio2.fi"
		expectedHostName1 = "127.0.0.1"
		expectedHostStr2  = "localhost local.test.fi local.test-app.fi"
		expectedHostName2 = "localhost"
	)

	// New domains struct.
	domains := models.NewSiteDomains()

	domains.SetDomain(&models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.luolamies.fi",
	})
	domains.SetDomain(&models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.amputaatio.fi",
	})
	domains.SetDomain(&models.Domain{
		Host:       "127.0.0.1",
		DomainName: "local.amputaatio2.fi",
	})
	domains.SetDomain(&models.Domain{
		Host:       "localhost",
		DomainName: "local.test.fi",
	})
	domains.SetDomain(&models.Domain{
		Host:       "localhost",
		DomainName: "local.test-app.fi",
	})

	// New hosts map.
	hostsMap := models.NewHostsMap()

	for _, domain := range domains.Domains {
		hostsMap.AddDomain(domain)
	}

	hostsDomains, err := hostsMap.GetHostDomains(expectedHostName1)
	if err != nil {
		t.Error(err)
	}

	if hostsDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", expectedHostName1, hostsDomains.String(), expectedHostStr1)
	}

	hostsDomains, err = hostsMap.GetHostDomains(expectedHostName2)
	if err != nil {
		t.Error(err)
	}

	if hostsDomains.String() != expectedHostStr2 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", expectedHostName2, hostsDomains.String(), expectedHostStr2)
	}
}
