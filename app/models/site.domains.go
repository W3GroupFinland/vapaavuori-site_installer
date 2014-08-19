package models

type SiteDomains struct {
	SiteName     string
	SubDirectory string
	Domains      map[string]*Domain
}

type Domain struct {
	Host       string
	DomainName string
}

func NewSiteDomains() *SiteDomains {
	return &SiteDomains{Domains: make(map[string]*Domain)}
}

func (sd *SiteDomains) SetDomain(d *Domain) {
	if d.DomainName == "" || d.Host == "" {
		return
	}

	if _, ok := sd.Domains[sd.getMapKey(d.DomainName, d.Host)]; !ok {
		sd.Domains[sd.getMapKey(d.DomainName, d.Host)] = d
	}
}

func (sd *SiteDomains) DomainExists(domain string, host string) bool {
	if _, ok := sd.Domains[sd.getMapKey(domain, host)]; ok {
		return true
	}

	return false
}

func (sd *SiteDomains) DeleteDomain(domain string, host string) {
	delete(sd.Domains, sd.getMapKey(domain, host))
}

func (sd *SiteDomains) getMapKey(domain string, host string) string {
	return domain + "_" + host
}
