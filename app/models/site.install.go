package models

type SiteInstallInfo struct {
	InstallRoot  string
	InstallType  string
	SiteName     string
	SubDirectory string
}

func NewSiteInstallInfo() *SiteInstallInfo {
	return &SiteInstallInfo{}
}
