package hosts_domains

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (h *HostsDomains) TestHostsDomainsParse(t *testing.T) {
	const (
		hostStr1         = "127.0.0.1 local.amputaatio.fi local.amputaatio2.fi local.luolamies.fi"
		expectedHostStr1 = "127.0.0.1 local.amputaatio.fi local.amputaatio2.fi local.ankkuri.fi local.luolamies.fi"
		host1            = "127.0.0.1"
		hostStr2         = "localhost local.test-app.fi local.test.fi"
		host2            = "localhost"
		hostStr3         = "focalhost "
	)

	// New hosts map.
	hostsMap := models.NewHostsMap()

	// New Host domains.
	hd := models.NewHostDomains()
	err := hd.Parse(hostStr1)
	if err != nil {
		t.Error(err)
	}
	hostsMap.AddHostDomains(hd)

	// Adding new domain to list.
	hostsMap.AddDomain(&models.Domain{
		Host:       host1,
		DomainName: "local.ankkuri.fi",
	})

	// New Host domains.
	hd = models.NewHostDomains()
	err = hd.Parse(hostStr2)
	if err != nil {
		t.Error(err)
	}
	hostsMap.AddHostDomains(hd)

	hostsDomains, err := hostsMap.GetHostDomains(host1)
	if err != nil {
		t.Error(err)
	}

	// New Host domains.
	hd = models.NewHostDomains()
	err = hd.Parse(hostStr3)
	if err == nil {
		t.Errorf("hd.Parse(hostStr3) = %v, expected = %v", err, "Not enough parts to parse from host string.")
	}
	hostsMap.AddHostDomains(hd)

	if hostsDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", host1, hostsDomains.String(), hostStr1)
	}

	hostsDomains, err = hostsMap.GetHostDomains(host2)
	if err != nil {
		t.Error(err)
	}

	if hostsDomains.String() != hostStr2 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", host2, hostsDomains.String(), hostStr2)
	}
}
