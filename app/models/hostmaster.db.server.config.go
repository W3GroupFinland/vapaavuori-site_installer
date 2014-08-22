package models

type DatabaseServerConfig struct {
	SiteId     int64
	ServerType string
	Template   string
	Port       int64
	ConfigRoot string
	ConfigFile string
	Secured    bool
	Cert       string
	CertKey    string
}
