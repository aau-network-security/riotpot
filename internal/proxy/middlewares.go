package proxy

import (
	"fmt"
	"net"
)

var (
	// Exportable middlewares manager
	Middlewares = NewMiddlewareManager()
)

// Use this interface to create new middlewares that include the `handle` function
type Middleware interface {
	// Handle a connection, do something to it
	handle(conn net.Conn) (net.Conn, error)
}

type MiddlewareManager interface {
	// Apply all the registered middlewares having the connection in consideration
	Apply(conn net.Conn) (ret net.Conn, err error)
	// Register a new middleware
	Register(middleware Middleware) (Middleware, error)
}

type MiddlewareManagerItem struct {
	MiddlewareManager

	// List of middlewares
	middlewares []Middleware
}

// Register a middleware
func (mm *MiddlewareManagerItem) Register(middleware Middleware) (mid Middleware, err error) {
	// Iterate the registered middlewares
	for _, md := range mm.middlewares {
		if md == middleware {
			err = fmt.Errorf("middleware already registered")
			return
		}
	}

	// Append the middleware to the list of registered middlewares
	mm.middlewares = append(mm.middlewares, middleware)
	mid = middleware

	return
}

// Apply each middleware to the connection
func (mm *MiddlewareManagerItem) Apply(conn net.Conn) (ret net.Conn, err error) {
	for _, middleware := range mm.middlewares {
		ret, err = middleware.handle(conn)
	}

	return
}

func NewMiddlewareManager() *MiddlewareManagerItem {
	return &MiddlewareManagerItem{
		// Create a slice of size 0 for the middlewares
		middlewares: make([]Middleware, 0),
	}
}
