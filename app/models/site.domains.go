package models

type SiteDomains struct {
	SiteName     string
	SubDirectory string
	Domains      map[string]string
}

func NewSiteDomains() *SiteDomains {
	return &SiteDomains{Domains: make(map[string]string)}
}

func (sd *SiteDomains) SetDomain(domain string) {
	if domain == "" {
		return
	}

	if _, ok := sd.Domains[domain]; !ok {
		sd.Domains[domain] = domain
	}
}

func (sd *SiteDomains) DomainExists(domain string) bool {
	if _, ok := sd.Domains[domain]; ok {
		return true
	}

	return false
}

func (sd *SiteDomains) DeleteDomain(domain string) {
	delete(sd.Domains, domain)
}
