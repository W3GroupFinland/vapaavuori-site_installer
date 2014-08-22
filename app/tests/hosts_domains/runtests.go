package hosts_domains

import (
	"testing"
)

type HostsDomains struct {
}

func Init() *HostsDomains {
	return &HostsDomains{}
}

func (hd *HostsDomains) RunTests(t *testing.T) {
	hd.TestHostsDomainsToString(t)
	hd.TestHostsDomainsParse(t)
	hd.TestHostsDomainsFromFile(t)
}
