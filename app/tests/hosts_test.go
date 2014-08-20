package tests

import (
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHostDomainsToString(t *testing.T) {
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

	const (
		expectedHostStr1  = "127.0.0.1 local.luolamies.fi local.amputaatio.fi local.amputaatio2.fi"
		expectedHostName1 = "127.0.0.1"
		expectedHostStr2  = "localhost local.test.fi local.test-app.fi"
		expectedHostName2 = "localhost"
	)

	hostDomains, err := hostsMap.GetHostDomains(expectedHostName1)
	if err != nil {
		t.Error(err)
	}

	if hostDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", expectedHostName1, hostDomains.String(), expectedHostStr1)
	}

	hostDomains, err = hostsMap.GetHostDomains(expectedHostName2)
	if err != nil {
		t.Error(err)
	}

	if hostDomains.String() != expectedHostStr2 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", expectedHostName2, hostDomains.String(), expectedHostStr2)
	}
}

func TestHostDomainsParse(t *testing.T) {
	const (
		hostStr1         = "127.0.0.1 local.luolamies.fi local.amputaatio.fi local.amputaatio2.fi"
		expectedHostStr1 = "127.0.0.1 local.luolamies.fi local.amputaatio.fi local.amputaatio2.fi local.ankkuri.fi"
		host1            = "127.0.0.1"
		hostStr2         = "localhost local.test.fi local.test-app.fi"
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

	hostDomains, err := hostsMap.GetHostDomains(host1)
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

	if hostDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", host1, hostDomains.String(), hostStr1)
	}

	hostDomains, err = hostsMap.GetHostDomains(host2)
	if err != nil {
		t.Error(err)
	}

	if hostDomains.String() != hostStr2 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", host2, hostDomains.String(), hostStr2)
	}
}

const (
	hostsContent = `##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1       localhost
255.255.255.255 broadcasthost
::1             localhost
fe80::1%lo0     localhost
#SITE_INSTALLER_HOSTS START
127.0.0.1 local.hogus.fi local.bogus.fi local.bim.fi
localhost local.exampleorg.fi local.exampleorg1.fi
#SITE_INSTALLER_HOSTS END`
)

func TestHostsDomainsFromFile(t *testing.T) {
	// Randomize temporary directory name.
	tempDir := "tmp-" + utils.RandomString(16)

	// Create temporary directory.
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Error(err)
	}

	const (
		expectedHostStr1  = "127.0.0.1 local.hogus.fi local.bogus.fi local.bim.fi local.example.org local.example-site.org"
		expectedHostName1 = "127.0.0.1"
		expectedHostStr2  = "localhost local.exampleorg.fi local.exampleorg1.fi local.ensemble.org"
		expectedHostName2 = "localhost"
	)

	// Write contents to file.
	fp := filepath.Join(tempDir, "hosts")
	err = ioutil.WriteFile(fp, []byte(hostsContent), 0644)
	if err != nil {
		t.Error(err)
	}

	site := controllers.Site{}
	hostsMap, err := site.ReadHostsFile(fp)
	if err != nil {
		t.Error(err)
	}
	hostsMap.AddDomain(&models.Domain{Host: expectedHostName1, DomainName: "local.example.org"})
	hostsMap.AddDomain(&models.Domain{Host: expectedHostName1, DomainName: "local.example-site.org"})
	hostsMap.AddDomain(&models.Domain{Host: expectedHostName2, DomainName: "local.ensemble.org"})
	hostsMap.AddDomain(&models.Domain{Host: "only.testing.net", DomainName: "local.ensemble.org"})

	hostDomains, err := hostsMap.GetHostDomains(expectedHostName1)
	if hostDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", expectedHostName1, hostDomains.String(), expectedHostStr1)
	}

	hostDomains, err = hostsMap.GetHostDomains(expectedHostName2)
	if hostDomains.String() != expectedHostStr2 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", expectedHostName2, hostDomains.String(), expectedHostStr2)
	}

	err = site.WriteNewHosts(fp, &hostsMap)
	if err != nil {
		t.Error(err)
	}

	const (
		extendedHostStr2 = expectedHostStr2 + " local.example.org"
	)

	mapAfter, err := site.ReadHostsFile(fp)
	// This should already exist in hosts map.
	mapAfter.AddDomain(&models.Domain{Host: expectedHostName1, DomainName: "local.example.org"})
	// This is added to hosts map.
	mapAfter.AddDomain(&models.Domain{Host: expectedHostName2, DomainName: "local.example.org"})

	hostDomains, err = mapAfter.GetHostDomains(expectedHostName1)
	if hostDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", expectedHostName1, hostDomains.String(), expectedHostStr1)
	}

	hostDomains, err = mapAfter.GetHostDomains(expectedHostName2)
	if hostDomains.String() != extendedHostStr2 {
		t.Errorf("hostsMap.GetHostDomains(\"%v\") = %v, expected = %v", expectedHostName2, hostDomains.String(), extendedHostStr2)
	}

	// Remove temporary directory.
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Error(err)
	}
}
