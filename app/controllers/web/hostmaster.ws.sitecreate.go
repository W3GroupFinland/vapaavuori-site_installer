package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	web_models "github.com/tuomasvapaavuori/site_installer/app/models/web"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func (c *HostmasterWS) InitInstallTemplate(platformId int64, tmpl *models.InstallTemplate) {
	tmpl.InstallInfo.PlatformId = platformId
	tmpl.InstallInfo.HttpUser = "_www"
	tmpl.InstallInfo.HttpGroup = "_www"
	tmpl.InstallInfo.DrupalRoot = filepath.Join(c.Base.Config.Platform.Directory, tmpl.InstallInfo.PlatformName)

	if tmpl.InstallInfo.TemplatePath != "" {
		tmpl.InstallInfo.TemplatePath = filepath.Join(c.Base.Config.SiteTemplates.Directory, tmpl.InstallInfo.TemplatePath)
	}

	tmpl.InstallInfo.InstallType = "template"

	// Set HTTP server values if template not empty.
	if tmpl.HttpServer.Include {
		tmpl.HttpServer.ConfigRoot = c.Base.Config.ServerConfigRoot.Http

		if tmpl.HttpServer.Template != "" {
			tmpl.HttpServer.Template = filepath.Join(c.Base.Config.SiteServerTemplates.Directory, tmpl.HttpServer.Template)
		}
	}

	// Set SSL server values if template not empty.
	if tmpl.SSLServer.Include {
		tmpl.SSLServer.Certificate = filepath.Join(c.Base.Config.SiteServerTemplates.Certificates, tmpl.SSLServer.Certificate)
		tmpl.SSLServer.ConfigRoot = c.Base.Config.ServerConfigRoot.SSL
		tmpl.SSLServer.Key = filepath.Join(c.Base.Config.SiteServerTemplates.Certificates, tmpl.SSLServer.Key)

		if tmpl.SSLServer.Template != "" {
			tmpl.SSLServer.Template = filepath.Join(c.Base.Config.SiteServerTemplates.Directory, tmpl.SSLServer.Template)
		}
	}
	tmpl.Init()

	// TODO: Make this nicer..
	tmpl.MysqlUser = models.RandomValue{Random: true}
	tmpl.MysqlPassword = models.RandomValue{Random: true}
	tmpl.DatabaseName = models.RandomValue{Random: true}
	tmpl.MysqlUserHosts = models.MysqlUserHosts{Hosts: []string{"127.0.0.1", "localhost"}}
	tmpl.MysqlUserPrivileges = models.MysqlUserPrivileges{Privileges: []string{"ALL"}}
	tmpl.MysqlGrantOption = models.MysqlGrantOption{Value: true}
}

func (c *HostmasterWS) ValidateInstallTemplate(tmpl *models.InstallTemplate) []*web_models.FormError {
	var errorList []*web_models.FormError

	exists, _, _ := c.HostMasterDB.PlatformExists(tmpl.InstallInfo.PlatformName, c.Base.Config.Platform.Directory)
	if !exists {
		errorList = append(errorList, web_models.NewFormError(
			"InstallInfo.PlatformName",
			"Does not exist.",
		))
	}

	if !utils.FileExists(tmpl.InstallInfo.DrupalRoot) {
		errorList = append(errorList, web_models.NewFormError(
			"InstallInfo.DrupalRoot",
			"Does not exist.",
		))
	}

	if tmpl.InstallInfo.TemplatePath == "" {
		errorList = append(errorList, web_models.NewFormError(
			"InstallInfo.TemplatePath",
			"Is empty",
		))
	}

	if !utils.FileExists(tmpl.InstallInfo.TemplatePath) {
		errorList = append(errorList, web_models.NewFormError(
			"InstallInfo.TemplatePath",
			"Does not exist.",
		))
	}

	// Http server values.
	if tmpl.HttpServer.Include {

		if tmpl.HttpServer.Template == "" {
			errorList = append(errorList, web_models.NewFormError(
				"HttpServer.Template",
				"Is empty",
			))
		}

		if !utils.FileExists(tmpl.HttpServer.Template) {
			errorList = append(errorList, web_models.NewFormError(
				"HttpServer.Template",
				"Does not exist.",
			))
		}

		if !utils.FileExists(tmpl.HttpServer.ConfigRoot) {
			errorList = append(errorList, web_models.NewFormError(
				"HttpServer.ConfigRoot",
				"Does not exist.",
			))
		}

		if exists, _ := c.HostMasterDB.DomainExists(tmpl.HttpServer.DomainInfo); exists {
			errorList = append(errorList, web_models.NewFormError(
				"HttpServer.DomainInfo",
				c.DomainExistsErrorStr(tmpl.HttpServer.DomainInfo),
			))
		}

		errorList = c.ValidateInstallTemplateDA(tmpl.HttpServer.DomainAliases, errorList)

	}

	// SSL server values.
	if tmpl.SSLServer.Include {

		if tmpl.SSLServer.Template == "" {
			errorList = append(errorList, web_models.NewFormError(
				"SSLServer.Template",
				"Is empty",
			))
		}

		if !utils.FileExists(tmpl.SSLServer.Template) {
			errorList = append(errorList, web_models.NewFormError(
				"SSLServer.Template",
				"Does not exist.",
			))
		}

		if !utils.FileExists(tmpl.SSLServer.ConfigRoot) {
			errorList = append(errorList, web_models.NewFormError(
				"SSLServer.ConfigRoot",
				"Does not exist.",
			))
		}

		if !utils.FileExists(tmpl.SSLServer.Certificate) {
			errorList = append(errorList, web_models.NewFormError(
				"SSLServer.Certificate",
				"Does not exist.",
			))
		}

		if !utils.FileExists(tmpl.SSLServer.Key) {
			errorList = append(errorList, web_models.NewFormError(
				"SSLServer.Key",
				"Does not exist.",
			))
		}

		if exists, _ := c.HostMasterDB.DomainExists(tmpl.SSLServer.DomainInfo); exists {
			errorList = append(errorList, web_models.NewFormError(
				"SSLServer.DomainInfo",
				c.DomainExistsErrorStr(tmpl.SSLServer.DomainInfo),
			))
		}

		errorList = c.ValidateInstallTemplateDA(tmpl.SSLServer.DomainAliases, errorList)

	}

	if exists, _ := c.HostMasterDB.SiteExists(tmpl.InstallInfo.PlatformId, tmpl.InstallInfo.SubDirectory); exists {
		errorList = append(errorList, web_models.NewFormError(
			"InstallInfo.SubDirectory",
			"Install subdirectory exists already.",
		))
	}

	if exists, _ := c.HostMasterDB.DomainExists(tmpl.InstallInfo.DomainInfo); exists {
		errorList = append(errorList, web_models.NewFormError(
			"InstallInfo.DomainInfo",
			c.DomainExistsErrorStr(tmpl.InstallInfo.DomainInfo),
		))
	}

	errorList = c.ValidateInstallTemplateDA(tmpl.InstallInfo.DomainAliases, errorList)

	return errorList
}

func (c *HostmasterWS) DomainExistsErrorStr(domain *models.Domain) string {
	return fmt.Sprintf("Domain %v exists already on host %v", domain.DomainName, domain.Host)
}

func (c *HostmasterWS) ValidateInstallTemplateDA(domains []*models.Domain, errorList []*web_models.FormError) []*web_models.FormError {
	for _, domain := range domains {
		if exists, _ := c.HostMasterDB.DomainExists(domain); exists {
			errorList = append(errorList, web_models.NewFormError(
				"Domain aliases",
				c.DomainExistsErrorStr(domain),
			))
		}
	}

	return errorList
}

func (c *HostmasterWS) RegisterFullSite(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	tmpl := models.NewInstallTemplate()
	// Get data from request to template model.
	c.RequestDataToModel(req, tmpl)

	exists, id, err := c.HostMasterDB.PlatformExists(tmpl.InstallInfo.PlatformName, c.Base.Config.Platform.Directory)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		log.Println(err)
		resp.SetError(http.StatusNotFound, "Platform doesn't exist.")
		return
	}

	// Initialize / correct values given to install template.
	c.InitInstallTemplate(id, tmpl)

	// Validate install template
	errorList := c.ValidateInstallTemplate(tmpl)
	if len(errorList) > 0 {
		resp.SetCallback(req).SetData(web_models.ResponseFormError, errorList)
		return
	}

	// Send ok after validation.
	resp.SetCallback(req).SetData(web_models.ResponseProcessStarting, nil)
	err = conn.WriteJSON(resp)
	if err != nil {
		log.Println("Websocket message error: %v", err.Error())
	}

	// Sleep for second to give browser time to prepare..
	time.Sleep(1 * time.Second)

	// Create process status messager.
	ps := web_models.NewWebProcess("Register full site", c.Channels.ConnStatusMsg, conn)
	ps.Start()

	sp, err := ps.AddSubProcess("Site Template Installation")
	if err != nil {
		log.Println(err)
	}

	_, err = c.Site.SiteTemplateInstallation(tmpl, sp)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	sp, err = ps.AddSubProcess("Create site")
	if err != nil {
		log.Println(err)
	}

	_, err = c.HostMasterDB.CreateSite(tmpl, sp)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	sp, err = ps.AddSubProcess("Write server configs")
	if err != nil {
		log.Println(err)
	}

	err = c.SiteTemplate.WriteApacheConfig(tmpl, sp)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	sp, err = ps.AddSubProcess("Write server configs to database")
	if err != nil {
		log.Println(err)
	}

	err = c.HostMasterDB.CreateServerConfigs(tmpl, sp)

	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	sp, err = ps.AddSubProcess("Create site symlinks")
	if err != nil {
		log.Println(err)
	}

	domains := c.Site.GetSiteTemplateDomains(tmpl)
	c.Site.CreateDomainSymlinks(tmpl, domains, sp)
	// Create site domains.

	sp, err = ps.AddSubProcess("Create site symlinks to database")
	if err != nil {
		log.Println(err)
	}

	err = c.HostMasterDB.CreateSiteDomains(tmpl, domains, sp)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	sp, err = ps.AddSubProcess("Add domains to hosts")
	if err != nil {
		log.Println(err)
	}

	err = c.Site.AddToHosts(tmpl, domains, sp)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	sp, err = ps.AddSubProcess("Restart http server")
	if err != nil {
		log.Println(err)
	}

	err = c.System.HttpServerRestart(sp)
	if err != nil {
		log.Println(err)
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	c.PlatformsUpdated()
	ps.Finish()

	resp.SetCallback(req).SetData(web_models.ResponseSiteCreated, nil)
}
