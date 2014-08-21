package models

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
	Type          string
	Template      string
	Port          string
	DomainInfo    *Domain
	DomainAliases []*Domain
	ConfigRoot    string `gcfg:"config-root"`
}

type SSLServerTemplate struct {
	HttpServerTemplate
	Certificate string
	Key         string
}

func (it *InstallTemplate) GetSiteInstallInfo() *SiteInstallInfo {
	return &it.InstallInfo
}
