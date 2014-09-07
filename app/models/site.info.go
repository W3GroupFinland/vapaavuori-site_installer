package models

import (
	"errors"
)

type SiteInfo struct {
	Id           int64
	PlatformId   int64
	Name         string
	DbName       string
	DbUser       string
	SubDirectory string
	InstallType  string
	Template     string
	Domains      []*Domain
	SSLConfigs   []*DatabaseServerConfig
	HTTPConfigs  []*DatabaseServerConfig
}

func (si *SiteInfo) AddDomains(domains ...*Domain) {
	for _, domain := range domains {
		si.Domains = append(si.Domains, domain)
	}
}

func (si *SiteInfo) AddConfigs(configs ...*DatabaseServerConfig) {
	for _, config := range configs {
		switch config.Secured {
		case true:
			si.SSLConfigs = append(si.SSLConfigs, config)
			break
		case false:
			si.HTTPConfigs = append(si.HTTPConfigs, config)
			break
		}
	}
}

// Get domain from Site info with domain name and host name.
func (si *SiteInfo) GetDomainByName(dn string, host string) (*Domain, error) {
	for _, domain := range si.Domains {
		if domain.DomainName == dn && domain.Host == host {
			return domain, nil
		}
	}

	return &Domain{}, errors.New("Domain doesn't exist in site.")
}

func (si *SiteInfo) MapDomainsToConfigs(domains []*Domain) {
	var configs = make(map[int64]*DatabaseServerConfig)

	// Populate server configs to one map.
	for _, config := range si.SSLConfigs {
		configs[config.Id] = config
	}
	for _, config := range si.HTTPConfigs {
		configs[config.Id] = config
	}

	for _, domain := range domains {
		if _, exists := configs[domain.ServerConfigId]; !exists {
			continue
		}

		configs[domain.ServerConfigId].AddDomain(domain)
	}
}
