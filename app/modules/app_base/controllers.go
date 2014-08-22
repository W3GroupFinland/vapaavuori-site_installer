package app_base

import (
	"log"
	"reflect"
)

type AppController interface {
	ControllerName() string
	Init()
}

type Controllers struct {
	Controllers map[string]*Controller
}

type Controller struct {
	AppController
}

// Controller definitions have to be translated to real controllers.
type ControllerDefinition struct {
	Name    string
	Handler string
}

func (c *Controllers) Get(name string) *Controller {
	if _, exists := c.Controllers[name]; !exists {
		log.Fatalf("Controller with name %v doesn't exist.", name)
	}

	return c.Controllers[name]
}

func (c *Controllers) Set(controller *Controller) {
	if c.Controllers == nil {
		c.Controllers = make(map[string]*Controller)
	}

	controller.Init()
	c.Controllers[controller.AppController.ControllerName()] = controller
}

func (c *Controller) Handler(name string) HttpHandlerFunc {
	_, exists := reflect.TypeOf(c.AppController).MethodByName(name)

	if !exists {
		log.Fatalf("Controller method %v doesn't exist.", name)
	}

	// Create function type convert to.
	var basicHandler HttpHandlerFunc
	basicHandlerType := reflect.TypeOf(basicHandler)

	// Get the value of the method name.
	methodVal := reflect.ValueOf(c.AppController).MethodByName(name)

	// Create converted function of handler.
	converted := methodVal.Convert(basicHandlerType)

	// Converted function to interface
	methodInterface := converted.Interface()
	// Turn method interface to corresponding http handler function.
	method := methodInterface.(HttpHandlerFunc)

	return method
}

func (c *Controller) AclHandler(name string) HttpAclHandlerFunc {
	_, exists := reflect.TypeOf(c.AppController).MethodByName(name)

	if !exists {
		log.Fatalf("Controller Acl method %v doesn't exist.", name)
	}

	// Create function type convert to.
	var basicHandler HttpAclHandlerFunc
	basicHandlerType := reflect.TypeOf(basicHandler)

	// Get the value of the method name.
	methodVal := reflect.ValueOf(c.AppController).MethodByName(name)

	// Create converted function of handler.
	converted := methodVal.Convert(basicHandlerType)

	// Converted function to interface
	methodInterface := converted.Interface()
	// Turn method interface to corresponding http handler function.
	method := methodInterface.(HttpAclHandlerFunc)

	return method
}

func NewController(app AppController) *Controller {
	c := new(Controller)
	c.AppController = app

	return c
}
