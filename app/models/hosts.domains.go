package models

import (
	"errors"
	"log"
	"strings"
)

type HostDomains struct {
	Host    string
	Domains map[string]*Domain
}

func NewHostDomains() *HostDomains {
	return &HostDomains{
		Domains: make(map[string]*Domain),
	}
}

func (hd *HostDomains) DomainExists(name string) bool {
	if _, ok := hd.Domains[name]; !ok {
		return false
	}

	return true
}

func (hd *HostDomains) AddDomain(domain *Domain) error {
	if hd.DomainExists(domain.DomainName) {
		log.Printf("Domain %v exists already on host %v.\n", domain.DomainName, hd.Host)
		return errors.New("Domain exists already.")
	}
	hd.Domains[domain.DomainName] = domain

	return nil
}

func (hd *HostDomains) RemoveDomain(name string) {
	delete(hd.Domains, name)
}

func (hd *HostDomains) Parse(str string) error {
	str = strings.TrimSpace(str)
	parts := strings.Split(str, " ")
	if len(parts) < 2 {
		return errors.New("Not enough parts to parse from host string.")
	}

	// First slice element is host name
	hd.Host = parts[0]
	// Rest of items are domains attached to host name.
	parts = parts[1:]

	for _, domainName := range parts {
		hd.AddDomain(&Domain{
			Host:       hd.Host,
			DomainName: domainName,
		})
	}

	return nil
}

func (hd *HostDomains) String() string {
	var items []string
	// First item of slice is host name
	items = append(items, hd.Host)

	// Append domain names to end of slice.
	for _, domain := range hd.Domains {
		items = append(items, domain.DomainName)
	}

	return strings.Join(items, " ")
}

func (hd *HostDomains) Bytes(sep byte) (out []byte) {
	out = append(out, []byte(hd.Host)...)
	out = append(out, sep)

	for _, domain := range hd.Domains {
		out = append(out, []byte(domain.DomainName)...)
		out = append(out, sep)
	}

	return out
}
