package controllers

import (
	"github.com/gorilla/websocket"
	web_models "github.com/tuomasvapaavuori/site_installer/app/models/web"
	"net/http"
)

func (c *HostmasterWS) GetServerTemplates(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	items, err := c.System.GetSiteServerTemplates()
	if err != nil {
		resp.SetError(http.StatusInternalServerError, err.Error())
		return
	}
	resp.SetCallback(req).SetData(web_models.ResponseServerTemplates, items)
}

func (c *HostmasterWS) GetSiteTemplates(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	items, err := c.System.GetSiteTemplates()
	if err != nil {
		resp.SetError(http.StatusInternalServerError, err.Error())
		return
	}
	resp.SetCallback(req).SetData(web_models.ResponseSiteTemplates, items)
}

func (c *HostmasterWS) GetServerCertificates(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	items, err := c.System.GetSiteServerCertificates()
	if err != nil {
		resp.SetError(http.StatusInternalServerError, err.Error())
		return
	}
	resp.SetCallback(req).SetData(web_models.ResponseServerCertificates, items)
}

func (c *HostmasterWS) GetUser(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	user, err := c.GetUserFromConnection(conn)
	if err != nil {
		resp.SetError(http.StatusBadRequest, "No user connection.")
		return
	}
	user.Password = "-- Obfuscated --"
	resp.SetCallback(req).SetData(web_models.ResponseUser, user)
}

func (c *HostmasterWS) GetPlatforms(conn *websocket.Conn, req *web_models.WebSocketRequest, resp *web_models.WebSocketResponse) {
	resp.SetCallback(req).SetData(web_models.ResponsePlatforms, c.Platforms.ToSliceList())
}
