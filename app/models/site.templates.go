package models

import (
	"errors"
)

type InstallTemplate struct {
	MysqlUser           RandomValue         `gcfg:"mysql-user"`
	MysqlPassword       RandomValue         `gcfg:"mysql-password"`
	MysqlUserHosts      MysqlUserHosts      `gcfg:"mysql-user-hosts"`
	MysqlUserPrivileges MysqlUserPrivileges `gcfg:"mysql-user-privileges"`
	MysqlGrantOption    MysqlGrantOption    `gcfg:"mysql-grant-option"`
	DatabaseName        RandomValue         `gcfg:"database-name"`
	InstallInfo         SiteInstallInfo     `gcfg:"install-info"`
	HttpServer          HttpServerTemplate  `gcfg:"http-server"`
	SSLServer           SSLServerTemplate   `gcfg:"ssl-server"`
	RollBack            *SiteRollBack
}

// Initialize install rollback.
func (it *InstallTemplate) Init() {
	it.RollBack = NewSiteRollBack(it)
}

type MysqlGrantOption struct {
	Value bool
}

type MysqlUserPrivileges struct {
	Privileges []string
}

type MysqlUserHosts struct {
	Hosts []string
}

type HttpServerTemplate struct {
	Type           string
	Template       string
	Port           int64
	DomainInfo     *Domain
	DomainAliases  []*Domain
	ConfigRoot     string `gcfg:"config-root"`
	ConfigFile     string
	ServerConfigId int64
}

type SSLServerTemplate struct {
	HttpServerTemplate
	Certificate string
	Key         string
}

// Helper function to update domain objects.
// Without site id and server id domain can't be created to database.
func (s *HttpServerTemplate) UpdateDomainIds(siteId int64) error {
	if s.ServerConfigId == 0 || siteId == 0 {
		return errors.New("Failed update domain ids. ServerConfigId or SiteId is zero.")
	}
	s.DomainInfo.ServerConfigId = s.ServerConfigId
	s.DomainInfo.SiteId = siteId
	for _, domain := range s.DomainAliases {
		domain.SiteId = siteId
		domain.ServerConfigId = s.ServerConfigId
	}

	return nil
}

func (it *InstallTemplate) GetSiteInstallInfo() *SiteInstallInfo {
	return &it.InstallInfo
}
