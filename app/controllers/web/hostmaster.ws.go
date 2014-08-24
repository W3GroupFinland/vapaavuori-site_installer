package controllers

import (
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
	Channels     struct {
		PlatformsUpdated  chan bool
		AllStatusMessages chan string
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
	// Start platforms updater routine.
	go c.PlatformsUpdater()
}

func (c *HostmasterWS) PlatformsUpdated() {
	c.Channels.PlatformsUpdated <- true
}

func (c *HostmasterWS) PlatformsUpdater() {
	for {
		updated := <-c.Channels.PlatformsUpdated
		if updated {
			log.Println("Platforms updated!")
			platforms := c.Platforms.ToSliceList()
			resp := web_models.WebSocketResponse{}
			resp.SetData(web_models.ResponsePlatforms, platforms).RefreshContent()

			for conn, _ := range c.Connections {
				err := conn.WriteJSON(resp)

				if err != nil {
					log.Println(err)
					continue
				}
			}
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
		resp.SetCallback(req).SetData(web_models.ResponsePlatforms, c.Platforms.ToSliceList())
		break

	// Register platform.
	case web_models.RequestRegisterPlatform:
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
		break

	// User info request.
	case web_models.RequestGetUser:
		user, err := c.GetUserFromConnection(conn)
		if err != nil {
			resp.SetError(http.StatusBadRequest, "No user connection.")
			return
		}
		user.Password = "-- Obfuscated --"
		resp.SetCallback(req).SetData(web_models.ResponseUser, user)
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
