package config

import (
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"strings"
)

type Command struct {
	Command   string
	Arguments []string
}

type Mysql struct {
	User     string
	Password string
	Protocol string
	Host     string
	Port     string
	DbName   string
}

type Host struct {
	Name string
	Port int
}

type Drush struct {
	Path string
}

type HttpServer struct {
	Restart string
}

type Backup struct {
	// TODO: Possible extra parameters to backup.
	Directory string
}

type Hosts struct {
	// TODO: Possible extra parameters to hosts.
	File string
}

type Platform struct {
	Directory string
}

type HttpSsl struct {
	HttpSsl     bool
	CertFile    string
	PrivateFile string
}

type SiteTemplates struct {
	Directory string
}

type SiteServerTemplates struct {
	Directory    string
	Certificates string
}

type ServerConfigRoot struct {
	Directory string
}

type Config struct {
	Host                Host
	Ssl                 HttpSsl
	Mysql               Mysql
	Drush               Drush
	HttpServer          HttpServer `gcfg:"http-server"`
	Backup              Backup
	Hosts               Hosts
	Platform            Platform
	SiteTemplates       SiteTemplates       `gcfg:"site-templates"`
	SiteServerTemplates SiteServerTemplates `gcfg:"site-server-templates"`
	ServerConfigRoot    ServerConfigRoot    `gcfg:"server-config-root"`
	WebTemplates        SiteTemplates       `gcfg:"web-templates"`
}

func NewConfig() *Config {
	return &Config{}
}

// Get configuration settings.
func (c *Config) Read(data []byte) {
	err := utils.ReadConfigData(data, c)
	if err != nil {
		log.Fatalln(err)
	}
}

func (hs *HttpServer) GetRestartCmd() (*Command, error) {
	// Parse restart command
	parts := strings.Split(hs.Restart, " ")

	if len(parts) == 0 {
		return &Command{}, errors.New("No http server restart command found.")
	}

	for idx, part := range parts {
		parts[idx] = strings.TrimSpace(part)
	}

	cmd := Command{}
	cmd.Command = parts[0]
	parts = parts[1:]

	for _, part := range parts {
		cmd.Arguments = append(cmd.Arguments, part)
	}

	return &cmd, nil
}
