package models

type InstallTemplate struct {
	MysqlUser           RandomValue         `gcfg:"mysql-user"`
	MysqlPassword       RandomValue         `gcfg:"mysql-password"`
	MysqlUserHosts      MysqlUserHosts      `gcfg:"mysql-user-hosts"`
	MysqlUserPrivileges MysqlUserPrivileges `gcfg:"mysql-user-privileges"`
	MysqlGrantOption    MysqlGrantOption    `gcfg:"mysql-grant-option"`
	DatabaseName        RandomValue         `gcfg:"database-name"`
	InstallInfo         SiteInstallConfig   `gcfg:"install-info"`
	HttpServer          HttpServerTemplate  `gcfg:"http-server"`
	SSLServer           SSLServerTemplate   `gcfg:"ssl-server"`
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
	ServerName    string   `gcfg:"server-name"`
	ServerAliases []string `gcfg:"server-alias"`
	ConfigRoot    string   `gcfg:"config-root"`
}

type SSLServerTemplate struct {
	HttpServerTemplate
	Certificate string
	Key         string
}

func (it *InstallTemplate) GetSiteInstallConfig() *SiteInstallConfig {
	return &it.InstallInfo
}
