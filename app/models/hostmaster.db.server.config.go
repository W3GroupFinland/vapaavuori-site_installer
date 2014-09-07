package models

import (
	"errors"
)

type DatabaseServerConfig struct {
	Id         int64
	SiteId     int64
	ServerType string
	Template   string
	Port       int64
	ConfigRoot string
	ConfigFile string
	Secured    bool
	Cert       string
	CertKey    string
	Domains    []*Domain
}

func (sc *DatabaseServerConfig) AddDomain(domain *Domain) {
	sc.Domains = append(sc.Domains, domain)
}

// Get domain from Server config with domain name and host name.
func (sc *DatabaseServerConfig) GetDomainByName(dn string, host string) (*Domain, error) {
	for _, domain := range sc.Domains {
		if domain.DomainName == dn && domain.Host == host {
			return domain, nil
		}
	}

	return &Domain{}, errors.New("Domain doesn't exist in server config.")
}
