package controllers

import (
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type SiteTemplate struct {
	Base *a.AppBase
}

func (c *SiteTemplate) ReadTemplate(file string) (*models.InstallTemplate, error) {
	templ := models.InstallTemplate{}

	err := utils.ReadConfigFile(file, &templ)
	if err != nil {
		return &templ, err
	}
	return &templ, nil
}

func (c *SiteTemplate) WriteApacheConfig(tmpl *models.InstallTemplate) error {
	outputFileName := tmpl.InstallInfo.SiteName

	if tmpl.HttpServer.Template != "" {
		// Write regular apache config.
		log.Println("Write regular apache config.")

		if tmpl.HttpServer.ServerName == "" {
			tmpl.HttpServer.ServerName = tmpl.InstallInfo.ServerName
		}

		if len(tmpl.HttpServer.ServerAliases) == 0 {
			tmpl.HttpServer.ServerAliases = tmpl.InstallInfo.ServerAliases
		}

		if tmpl.HttpServer.ConfigRoot == "" {
			tmpl.HttpServer.ConfigRoot = tmpl.InstallInfo.ServerConfigRoot
		}

		log.Println("http." + outputFileName)

		err := c.WriteServerConfig(tmpl, "http."+outputFileName, tmpl.HttpServer.Template, tmpl.HttpServer.ConfigRoot)
		if err != nil {
			return err
		}
	}
	if tmpl.SSLServer.Template != "" {
		// Write SSL apache config.
		log.Println("Write SSL config.")

		if tmpl.SSLServer.ServerName == "" {
			tmpl.SSLServer.ServerName = tmpl.InstallInfo.ServerName
		}

		if len(tmpl.SSLServer.ServerAliases) == 0 {
			tmpl.SSLServer.ServerAliases = tmpl.InstallInfo.ServerAliases
		}

		if tmpl.SSLServer.ConfigRoot == "" {
			tmpl.SSLServer.ConfigRoot = tmpl.InstallInfo.ServerConfigRoot
		}

		log.Println("ssl." + outputFileName)

		err := c.WriteServerConfig(tmpl, "ssl."+outputFileName, tmpl.SSLServer.Template, tmpl.SSLServer.ConfigRoot)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *SiteTemplate) WriteServerConfig(tmpl *models.InstallTemplate, outputFileName string, templateFile string, configRoot string) error {
	t := template.Must(template.New("server.config").ParseFiles(templateFile))

	fo, err := os.Create(filepath.Join(configRoot, outputFileName))
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	err = t.ExecuteTemplate(fo, "server.config", tmpl)
	if err != nil {
		return err
	}

	return nil
}
