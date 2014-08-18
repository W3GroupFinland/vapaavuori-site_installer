package models

type SiteInstallConfig struct {
	DrupalRoot       string   `gcfg:"drupal-root"`
	InstallType      string   `gcfg:"install-type"`
	TemplatePath     string   `gcfg:"template-path"`
	SiteName         string   `gcfg:"sitename"`
	ServerName       string   `gcfg:"server-name"`
	ServerAliases    []string `gcfg:"server-alias"`
	SubDirectory     string   `gcfg:"sub-directory"`
	HttpUser         string   `gcfg:"http-user"`
	HttpGroup        string   `gcfg:"http-group"`
	ServerConfigRoot string   `gcfg:"server-config-root"`
}

func NewSiteInstallConfig() *SiteInstallConfig {
	return &SiteInstallConfig{}
}

type SiteRootInfo struct {
	DrupalVersion       string
	DefaultTheme        string
	AdministrationTheme string
	PHPConfig           string
	PHPOs               string
	DrushVersion        string
	DrushConfiguration  string
	DrushAliasFiles     string
	DrupalRoot          string
}
