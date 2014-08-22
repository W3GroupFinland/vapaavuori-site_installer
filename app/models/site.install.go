package models

type SiteInstallInfo struct {
	DrupalRoot       string `gcfg:"drupal-root"`
	InstallType      string `gcfg:"install-type"`
	TemplatePath     string `gcfg:"template-path"`
	SiteName         string `gcfg:"sitename"`
	DomainInfo       *Domain
	DomainAliases    []*Domain
	SubDirectory     string `gcfg:"sub-directory"`
	HttpUser         string `gcfg:"http-user"`
	HttpGroup        string `gcfg:"http-group"`
	ServerConfigRoot string `gcfg:"server-config-root"`
	PlatformId       int64
	PlatformName     string
	SiteId           int64
}

func NewSiteInstallConfig() *SiteInstallInfo {
	return &SiteInstallInfo{}
}
