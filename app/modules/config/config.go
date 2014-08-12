package config

import (
	"bitbucket.org/kardianos/osext"
	"code.google.com/p/gcfg"
	"log"
	"os"
	"path/filepath"
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

type Config struct {
	Host  Host
	Mysql Mysql
}

func NewConfig() *Config {
	return &Config{}
}

// Get Servicesuration settings from settings.gcfg
func (c *Config) Read(file string) {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(folderPath, file)
	absPath, err := GetAbsDirectory(path)
	if err != nil {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}

		absPath, err = GetAbsDirectory(dir)
		if err != nil {
			log.Fatalln(err)
		}

		path = filepath.Join(absPath, file)
	}

	err = gcfg.ReadFileInto(c, path)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAbsDirectory(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
