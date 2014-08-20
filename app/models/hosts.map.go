package models

import (
	"bytes"
	"errors"
	"strings"
	//"log"
)

// Set local constants to indicate hosts read state.
const (
	READ_NOT_STARTED = 0
	READ_STARTED     = 1
	READ_ENDED       = 2

	HOSTS_START_READ_STR = "#SITE_INSTALLER_HOSTS START"
	HOSTS_END_READ_STR   = "#SITE_INSTALLER_HOSTS END"
	NEW_LINE_BYTE        = 10
	SPACE_BYTE           = 32
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

func (hm HostsMap) AddDomain(d *Domain) error {
	if !hm.HostExists(d.Host) {
		hd := NewHostDomains()
		hd.Host = d.Host
		hm[d.Host] = hd
	}

	err := hm[d.Host].AddDomain(d)

	return err
}

func (hm HostsMap) String(sep string) (out string) {
	var parts []string
	// Snippet starting point.
	parts = append(parts, HOSTS_START_READ_STR)

	for _, host := range hm {
		parts = append(parts, host.String())
	}

	// Snippet ending point.
	parts = append(parts, HOSTS_END_READ_STR)

	return strings.Join(parts, sep)
}

func (hm HostsMap) Bytes(strSep byte, lineSep byte) (out []byte) {
	var parts [][]byte

	// Append start of application hosts string
	parts = append(parts, []byte(HOSTS_START_READ_STR))

	// Append every host name from host map.
	for _, host := range hm {
		parts = append(parts, host.Bytes(strSep))
	}

	// Append end of application hosts string
	parts = append(parts, []byte(HOSTS_END_READ_STR))
	out = bytes.Join(parts, []byte{NEW_LINE_BYTE})

	return out
}
