package models

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

type PlatformInfo struct {
	Name       string
	Registered bool
	PlatformId int64
	RootInfo   *SiteRootInfo
}
