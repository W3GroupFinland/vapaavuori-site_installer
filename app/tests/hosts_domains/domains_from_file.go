package hosts_domains

import (
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func (h *HostsDomains) TestHostsDomainsFromFile(t *testing.T) {
	const (
		expectedHostStr1  = "127.0.0.1 local.bim.fi local.bogus.fi local.example-site.org local.example.org local.hogus.fi"
		expectedHostName1 = "127.0.0.1"
		expectedHostStr2  = "localhost local.ensemble.org local.exampleorg.fi local.exampleorg1.fi"
		expectedHostName2 = "localhost"
	)

	// Randomize temporary directory name.
	tempDir := "tmp-" + utils.RandomString(16)

	// Create temporary directory.
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Error(err)
	}

	// Write contents to file.
	fp := filepath.Join(tempDir, "hosts")
	err = ioutil.WriteFile(fp, []byte(hostsContent), 0644)
	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadFile(fp)
	t.Log(string(b))

	site := controllers.Site{}
	hostsMap, err := site.ReadHostsFile(fp)
	if err != nil {
		t.Error(err)
	}
	hostsDomains, err := hostsMap.GetHostDomains(expectedHostName1)
	t.Log(hostsDomains.String())

	hostsMap.AddDomain(&models.Domain{Host: expectedHostName1, DomainName: "local.example.org"})
	hostsMap.AddDomain(&models.Domain{Host: expectedHostName1, DomainName: "local.example-site.org"})
	hostsMap.AddDomain(&models.Domain{Host: expectedHostName2, DomainName: "local.ensemble.org"})
	hostsMap.AddDomain(&models.Domain{Host: "only.testing.net", DomainName: "local.ensemble.org"})

	hostsDomains, err = hostsMap.GetHostDomains(expectedHostName1)
	if hostsDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", expectedHostName1, hostsDomains.String(), expectedHostStr1)
	}

	hostsDomains, err = hostsMap.GetHostDomains(expectedHostName2)
	if hostsDomains.String() != expectedHostStr2 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", expectedHostName2, hostsDomains.String(), expectedHostStr2)
	}

	err = site.WriteNewHosts(fp, &hostsMap)
	if err != nil {
		t.Error(err)
	}

	const (
		extendedHostStr2 = expectedHostStr2 + " local.xanus.org"
	)

	mapAfter, err := site.ReadHostsFile(fp)
	// This should already exist in hosts map.
	mapAfter.AddDomain(&models.Domain{Host: expectedHostName1, DomainName: "local.example.org"})
	// This is added to hosts map.
	mapAfter.AddDomain(&models.Domain{Host: expectedHostName2, DomainName: "local.xanus.org"})

	hostsDomains, err = mapAfter.GetHostDomains(expectedHostName1)
	if hostsDomains.String() != expectedHostStr1 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", expectedHostName1, hostsDomains.String(), expectedHostStr1)
	}

	hostsDomains, err = mapAfter.GetHostDomains(expectedHostName2)
	if hostsDomains.String() != extendedHostStr2 {
		t.Errorf("hostsMap.GetHostsDomains(\"%v\") = %v, expected = %v", expectedHostName2, hostsDomains.String(), extendedHostStr2)
	}

	// Remove temporary directory.
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Error(err)
	}
}
