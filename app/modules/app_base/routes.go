package app_base

import (
	"fmt"
	"log"
	"net/http"
)

type HttpHandlerFunc func(rw http.ResponseWriter, r *http.Request)
type HttpAclHandlerFunc func(rw http.ResponseWriter, r *http.Request) bool

type Route struct {
	Path             string
	Handlers         []RouteHandler
	HandlersByMethod map[string]RouteHandler
	MiddleWare       []HttpHandlerFunc
	Acl              []HttpAclHandlerFunc
}

type RouteHandler struct {
	Name       string
	Method     string
	Handler    HttpHandlerFunc
	MiddleWare []HttpHandlerFunc
}

// Route definition makes it easier to read routes
// from text file. Routes have to be translated to real routes afterwards.
// Route Controller handlers are reflect value functions.
type RouteDefinition struct {
	Path       string
	Handlers   []RouteDefinitionHandler
	MiddleWare []ControllerDefinition
}

type RouteDefinitionHandler struct {
	Method     string
	Handler    ControllerDefinition
	MiddleWare []ControllerDefinition
}

func (a *AppBase) AddRoute(r *Route) *Route {
	if a.Routes == nil {
		a.Routes = make(map[string]*Route)
	}

	if _, exists := a.Routes[r.Path]; exists {
		log.Fatalf("Trying to reregister path %v", r.Path)
	}

	if r.HandlersByMethod == nil {
		r.HandlersByMethod = make(map[string]RouteHandler)
	}

	for _, handler := range r.Handlers {
		r.HandlersByMethod[handler.Method] = handler
	}

	a.Routes[r.Path] = r

	return r
}

func (r *Route) AddMiddleWare(fn HttpHandlerFunc) *Route {
	r.MiddleWare = append(r.MiddleWare, fn)

	return r
}

func (r *Route) Handle() (string, HttpHandlerFunc) {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("Used method was: %v\n", req.Method)

		var (
			exists  bool
			handler RouteHandler
		)

		// Check that given mehtod for route exists.
		if handler, exists = r.HandlersByMethod[req.Method]; !exists {
			http.Error(rw, HttpErrorMsg(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Run current route ACL-functions
		for _, aclFunc := range r.Acl {
			ok := aclFunc(rw, req)
			if !ok {
				http.Error(rw, HttpErrorMsg(http.StatusForbidden), http.StatusForbidden)
				return
			}
		}

		// Run current route Middleware functions
		for _, mwFunc := range r.MiddleWare {
			mwFunc(rw, req)
		}

		// Run main controller function for route and method.
		handler.Handler(rw, req)
	}

	return r.Path, fn
}

func (a *AppBase) RegisterRoutes() {
	for _, route := range a.Routes {
		fmt.Printf("Registered route %v\n", route.Path)
		http.HandleFunc(route.Handle())
	}
}

func (a *AppBase) RoutesFromDefinitions(defs []RouteDefinition) {
	routes := a.TranslateRoutes(defs)

	for _, route := range routes {
		a.AddRoute(route)
	}
}

// Translates route definitions to application routes.
func (a *AppBase) TranslateRoutes(defs []RouteDefinition) map[string]*Route {
	routes := make(map[string]*Route)

	for _, def := range defs {
		var mw []HttpHandlerFunc
		for _, mwHandler := range def.MiddleWare {
			mw = append(mw, a.WebControllers.Get(mwHandler.Name).Handler(mwHandler.Handler))
		}

		routes[def.Path] = &Route{
			Path:       def.Path,
			MiddleWare: mw}

		for _, hd := range def.Handlers {
			routes[def.Path].Handlers = append(routes[def.Path].Handlers, RouteHandler{
				Method:  hd.Method,
				Handler: a.WebControllers.Get(hd.Handler.Name).Handler(hd.Handler.Handler)})
		}
	}

	return routes
}
