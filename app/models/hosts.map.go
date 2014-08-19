package models

import (
	"errors"
	//"log"
)

type HostsMap map[string]*HostDomains

func NewHostsMap() HostsMap {
	hm := HostsMap{}
	hm = make(map[string]*HostDomains)

	return hm
}

func (hm HostsMap) AddHostDomains(hd *HostDomains) {
	for _, domain := range hd.Domains {
		hm.AddDomain(domain)
	}
}

func (hm HostsMap) HostExists(host string) bool {
	if _, ok := hm[host]; !ok {
		return false
	}

	return true
}

func (hm HostsMap) GetHostDomains(host string) (*HostDomains, error) {
	if !hm.HostExists(host) {
		return NewHostDomains(), errors.New("Host with domains doesn't exist in map.")
	}

	return hm[host], nil
}

func (hm HostsMap) AddDomain(d *Domain) {
	if !hm.HostExists(d.Host) {
		hd := NewHostDomains()
		hd.Host = d.Host
		hm[d.Host] = hd
	}

	hm[d.Host].AddDomain(d)
}
