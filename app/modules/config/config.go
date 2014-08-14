package config

import (
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
)

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

type Config struct {
	Host  Host
	Mysql Mysql
	Drush Drush
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
