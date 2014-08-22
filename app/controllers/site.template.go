package controllers

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	a "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
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

		if tmpl.HttpServer.DomainInfo.DomainName == "" {
			tmpl.HttpServer.DomainInfo = tmpl.InstallInfo.DomainInfo
		}

		if len(tmpl.HttpServer.DomainAliases) == 0 {
			tmpl.HttpServer.DomainAliases = tmpl.InstallInfo.DomainAliases
		}

		if tmpl.HttpServer.ConfigRoot == "" {
			tmpl.HttpServer.ConfigRoot = tmpl.InstallInfo.ServerConfigRoot
		}

		outFile := "http." + outputFileName
		log.Println(outFile)
		err := c.WriteServerConfig(tmpl, outFile, tmpl.HttpServer.Template, tmpl.HttpServer.ConfigRoot)
		if err != nil {
			return err
		}

		tmpl.HttpServer.ConfigFile = outFile
	}
	if tmpl.SSLServer.Template != "" {
		// Write SSL apache config.
		log.Println("Write SSL config.")

		if tmpl.SSLServer.DomainInfo.DomainName == "" {
			tmpl.SSLServer.DomainInfo = tmpl.InstallInfo.DomainInfo
		}

		if len(tmpl.SSLServer.DomainAliases) == 0 {
			tmpl.SSLServer.DomainAliases = tmpl.InstallInfo.DomainAliases
		}

		if tmpl.SSLServer.ConfigRoot == "" {
			tmpl.SSLServer.ConfigRoot = tmpl.InstallInfo.ServerConfigRoot
		}

		outFile := "ssl." + outputFileName
		log.Println(outFile)
		err := c.WriteServerConfig(tmpl, outFile, tmpl.SSLServer.Template, tmpl.SSLServer.ConfigRoot)
		if err != nil {
			return err
		}

		tmpl.SSLServer.ConfigFile = outFile
	}

	return nil
}

func (c *SiteTemplate) WriteServerConfig(tmpl *models.InstallTemplate, outputFileName string, templateFile string, configRoot string) error {
	t := template.Must(template.New("server.config").ParseFiles(templateFile))

	fullPath := filepath.Join(configRoot, outputFileName)
	fo, err := os.Create(fullPath)
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

	// Rollback: Remove server config file.
	tmpl.RollBack.AddFileFunction(utils.RemoveFile, fullPath)

	return nil
}
