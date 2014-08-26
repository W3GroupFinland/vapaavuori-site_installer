package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	web_models "github.com/tuomasvapaavuori/site_installer/app/models/web"
	"net/http"
)

func (c *HostmasterWS) RegisterPlatform(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	pi := &models.PlatformInputRequest{}
	c.MapToStruct(pi, req.Data)

	// Get platform info from prepopulated data.
	platform, err := c.Platforms.Get(c.Base.Config.Platform.Directory, pi.Name)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Error())
		return
	}

	// Create template of platform register data.
	tmpl := models.InstallTemplate{
		InstallInfo: models.SiteInstallInfo{
			PlatformName: platform.Name,
		}}

	// Initialize rollback functionality.
	tmpl.Init()
	// Create platfrom according the template.
	id, err := c.HostMasterDB.CreatePlatform(&tmpl)
	if err != nil {
		resp.SetError(http.StatusInternalServerError, err.Error())
		tmpl.RollBack.Execute()
		return
	}

	// Set new values to object.
	platform.Registered = true
	platform.PlatformId = id

	// Set response data of successfull registration.
	resp.SetCallback(req).SetData(web_models.ResponsePlatformRegistered, platform)
	c.PlatformsUpdated()
	msg := fmt.Sprintf("Platform %v was updated.", platform.Name)
	c.StatusMessageToAll(msg)
}
