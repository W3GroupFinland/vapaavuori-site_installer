package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/tuomasvapaavuori/site_installer/app/controllers"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	web_models "github.com/tuomasvapaavuori/site_installer/app/models/web"
	a "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"log"
	"net/http"
	"reflect"
)

type HostmasterWS struct {
	Base         *a.AppBase
	User         *User // Pointer to user controller.
	System       *controllers.System
	Connections  map[*websocket.Conn]*models.User
	Upgrader     *websocket.Upgrader
	HostMasterDB *controllers.HostMasterDB
	Platforms    models.PlatformList
	Site         *controllers.Site
	SiteTemplate *controllers.SiteTemplate
	Channels     struct {
		PlatformsUpdated chan bool
		AllStatusMsg     chan *web_models.StatusMessage
		ConnStatusMsg    chan *web_models.ConnStatusMessage
	}
}

func (c *HostmasterWS) ControllerName() string {
	return "app.controllers.web.hostmaster.ws"
}

func (c *HostmasterWS) Init() {
	c.Connections = make(map[*websocket.Conn]*models.User)
	c.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Populate initial platform list.
	platforms, err := c.System.GetDrupalPlatforms()
	if err != nil {
		log.Fatal("Error getting drupal platforms in initialization.")
	}
	c.Platforms = platforms

	// Initialize channels
	c.Channels.PlatformsUpdated = make(chan bool)
	c.Channels.AllStatusMsg = make(chan *web_models.StatusMessage)
	c.Channels.ConnStatusMsg = make(chan *web_models.ConnStatusMessage)
	// Start platforms updater routine.
	go c.StatusUpdater()
}

func (c *HostmasterWS) PlatformsUpdated() {
	c.Channels.PlatformsUpdated <- true
}

func (c *HostmasterWS) MessageToConnection(conn *websocket.Conn, msg string) {
	data := web_models.NewConnStatusMessage(conn, msg)
	c.Channels.ConnStatusMsg <- data
}

func (c *HostmasterWS) StatusMessageToAll(msg string) {
	data := web_models.NewStatusMessage(msg)
	c.Channels.AllStatusMsg <- data
}

func (c *HostmasterWS) StatusUpdater() {
	for {
		resp := web_models.WebSocketResponse{}
		select {

		// Platforms updated.
		case updated := <-c.Channels.PlatformsUpdated:
			if updated {
				log.Println("Platforms updated!")
				platforms := c.Platforms.ToSliceList()
				resp.SetData(web_models.ResponsePlatforms, platforms).RefreshContent()
				c.WebSocketSendAll(resp)
			}

		// Status message to all connections.
		case msg := <-c.Channels.AllStatusMsg:
			resp.SetData(web_models.ResponseStatusMessage, msg).RefreshContent()
			c.WebSocketSendAll(resp)

		// Status message to one connection.
		case msg := <-c.Channels.ConnStatusMsg:
			// Get connection.
			conn := msg.Connection
			// Set new message.
			data := web_models.StatusMessage{Message: msg.Message}
			resp.SetData(web_models.ResponseStatusMessage, data).RefreshContent()
			err := conn.WriteJSON(resp)
			if err != nil {
				c.ConnectionDelete(conn)
				log.Println(err)
				continue
			}
		}
	}
}

func (c *HostmasterWS) WebSocketSendAll(i interface{}) {
	for conn, _ := range c.Connections {
		err := conn.WriteJSON(i)

		if err != nil {
			c.ConnectionDelete(conn)
			log.Println(err)
			continue
		}
	}
}

func (c *HostmasterWS) Messager(rw http.ResponseWriter, r *http.Request) {
	conn, err := c.Upgrader.Upgrade(rw, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(rw, "Not a websocket handshake.", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	user, ok := c.User.Current(rw, r)
	if !ok {
		http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
	}

	// Delete and close connection when finished.
	defer c.ConnectionDelete(conn)
	// Adds connection to controller connection map.
	c.AddConnection(conn, user)

	for {
		request := web_models.WebSocketRequest{}
		err := conn.ReadJSON(&request)
		if err != nil {
			log.Printf("Websocket message error: %v", err.Error())
			// Returns and closes connection.
			return
		}
		resp := c.RequestResolver(conn, &request)

		err = conn.WriteJSON(resp)
		if err != nil {
			log.Println("Websocket message error: %v", err.Error())
			return
		}
	}
}

func (c *HostmasterWS) AddConnection(conn *websocket.Conn, user *models.User) {
	c.Connections[conn] = user
}

func (c *HostmasterWS) ConnectionDelete(conn *websocket.Conn) {
	delete(c.Connections, conn)
	conn.Close()
}

func (c *HostmasterWS) GetUserFromConnection(conn *websocket.Conn) (*models.User, error) {
	if user, ok := c.Connections[conn]; ok {
		return user, nil
	}

	return &models.User{}, errors.New("Connection doesn't exist.")
}

// TODO: Refactor to smaller tasks.
func (c *HostmasterWS) RequestResolver(conn *websocket.Conn, req *web_models.WebSocketRequest) (resp *web_models.WebSocketResponse) {
	resp = &web_models.WebSocketResponse{}

	switch req.Type {
	// Platform request.
	case web_models.RequestGetPlatforms:
		c.GetPlatforms(conn, req, resp)
		break

	// Register platform.
	case web_models.RequestRegisterPlatform:
		c.RegisterPlatform(conn, req, resp)
		break

	case web_models.RequestRegisterFullSite:
		c.RegisterFullSite(conn, req, resp)
		break

	case web_models.RequestGetSiteTemplates:
		c.GetSiteTemplates(conn, req, resp)
		break

	case web_models.RequestGetServerTemplates:
		c.GetServerTemplates(conn, req, resp)
		break

	case web_models.RequestGetServerCertificates:
		c.GetServerCertificates(conn, req, resp)
		break

	// User info request.
	case web_models.RequestGetUser:
		c.GetUser(conn, req, resp)
		break
	default:
		resp.SetError(http.StatusBadRequest, "Bad request.")
	}

	return
}

func (c *HostmasterWS) MapToStruct(targetStruct interface{}, data interface{}) {
	values := data.(map[string]interface{})
	target := reflect.ValueOf(targetStruct).Elem()

	for idx, val := range values {
		if val == nil {
			continue
		}
		valTyp := reflect.TypeOf(val)
		field := target.FieldByName(idx)
		// Check that given field exists in struct.
		if !field.IsValid() {
			log.Printf("Invalid field %v given on WebSocketPost", idx)
			continue
		}
		switch field.Kind() {
		case reflect.String:
			if valTyp.Kind() != reflect.String {
				continue
			}
			field.SetString(val.(string))
			break
		case reflect.Int:
			if valTyp.Kind() != reflect.Float64 {
				continue
			}
			f := val.(float64)
			field.SetInt(int64(f))
			break
		case reflect.Int64:
			if valTyp.Kind() != reflect.Float64 {
				continue
			}
			f := val.(float64)
			field.SetInt(int64(f))
			break
		}
	}
}

func (c *HostmasterWS) RequestDataToModel(req *web_models.WebSocketRequest, i interface{}) error {
	data, err := json.Marshal(req.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}

	return nil
}
