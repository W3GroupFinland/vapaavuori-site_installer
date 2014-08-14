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
	configName := tmpl.InstallInfo.SiteName
	t := template.Must(template.New("server.config").ParseFiles(tmpl.HttpServer.Template))

	fo, err := os.Create(filepath.Join(tmpl.HttpServer.ConfigRoot, configName))
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
