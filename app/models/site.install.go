package models

type SiteInstallConfig struct {
	DrupalRoot   string
	InstallType  string
	SiteName     string
	SubDirectory string
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
