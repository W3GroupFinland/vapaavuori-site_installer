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
	Port string
}

type Drush struct {
	Path string
}

type HttpServer struct {
	Restart string
}

type Backup struct {
	Directory string
}

type Config struct {
	Host       Host
	Mysql      Mysql
	Drush      Drush
	HttpServer HttpServer `gcfg:"http-server"`
	Backup     Backup
}

func NewConfig() *Config {
	return &Config{}
}

// Get configuration settings from file.
func (c *Config) Read(file string) {
	err := utils.ReadConfigFile(file, c)
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
