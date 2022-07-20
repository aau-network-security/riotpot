package api

import (
	"github.com/gin-gonic/gin"
)

// Router interface
type Router interface {
	// Returns the list of routes registered in the router
	Routes() []Route
	// Parent path for the router
	Path() string
	// Group the registered path
	AddToGroup(parentGroup *gin.RouterGroup) *gin.RouterGroup
}

type AbstractRouter struct {
	// Inherits the Router
	Router
	// Routes registered in the router
	routes []Route
	// Parent path
	path string
	// Child Routers
	// A router may have child routers
	// For example:
	// - /v1
	// .. - /router1
	// .. .. - /r1
	// .. .. - /r2
	// .. - /router2
	childs []Router
}

// Returns the parent path
func (r *AbstractRouter) Path() string {
	return r.path
}

// Returns the list of routes registered in the router
func (r *AbstractRouter) Routes() []Route {
	return r.routes
}

// Add handlers to a router group
func (r *AbstractRouter) addHandlers(parentGroup *gin.RouterGroup) *gin.RouterGroup {
	// Iterate the routes and add the handlers registered in the
	for _, route := range r.Routes() {
		parentGroup.Handle(route.Method(), route.Path(), route.Handlers()...)
	}

	return parentGroup
}

// Add all the child routers to a group
func (r *AbstractRouter) addChilds(parentGroup *gin.RouterGroup) *gin.RouterGroup {
	// Iterate through the child routers to add the routes
	if len(r.childs) > 0 {
		for _, child := range r.childs {
			child.AddToGroup(parentGroup)
		}
	}

	return parentGroup
}

// Create the group routes inside of the router
func (r *AbstractRouter) AddToGroup(parentGroup *gin.RouterGroup) *gin.RouterGroup {
	// Create a group inside of the parent group for this child
	childGroup := parentGroup.Group(r.Path())
	// Add the routes handlers for the current group
	childGroup = r.addHandlers(childGroup)
	// Add the child groups to the router
	childGroup = r.addChilds(childGroup)

	return childGroup
}

func NewRouter(path string, routes []Route, childs []Router) *AbstractRouter {
	return &AbstractRouter{
		path:   path,
		routes: routes,
		childs: childs,
	}
}

// Route that will be handled by the API
type Route interface {
	// Raw function to handle the request
	Handlers() gin.HandlersChain
	// (Sub)Path to the route
	Path() string
	// Method used for the path
	Method() string
}

type AbstractRoute struct {
	Route
	path     string
	method   string
	handlers gin.HandlersChain
}

func (ar *AbstractRoute) Path() string {
	return ar.path
}

func (ar *AbstractRoute) Method() string {
	return ar.method
}

func (ar *AbstractRoute) Handlers() gin.HandlersChain {
	return ar.handlers
}

func NewRoute(path string, method string, handlers ...gin.HandlerFunc) *AbstractRoute {
	return &AbstractRoute{
		path:     path,
		method:   method,
		handlers: handlers,
	}
}
