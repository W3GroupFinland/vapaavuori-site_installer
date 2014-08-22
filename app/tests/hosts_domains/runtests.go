package hosts_domains

import (
	"testing"
)

type HostsDomains struct {
}

func RunTests(t *testing.T) {
	hd := HostsDomains{}
	hd.TestHostsDomainsToString(t)
	hd.TestHostsDomainsParse(t)
	hd.TestHostsDomainsFromFile(t)
}
