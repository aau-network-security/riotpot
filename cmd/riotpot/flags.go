/*
This package determines the flags set in the environment
*/

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"github.com/riotpot/api"
	"github.com/riotpot/api/proxy"
	"github.com/riotpot/api/service"
	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/rs/zerolog"

	_ "github.com/riotpot/statik"
)

type Routers []api.Router

var (
	// Routers of the application
	routers Routers = Routers{
		// Proxy router
		proxy.ProxiesRouter,
		// Services router
		service.ServicesRouter,
	}
)

var (
	debug  = flag.Bool("debug", true, "Set log level to debug")
	runApi = flag.Bool("api", true, "Whether to start the API")
)

func setupApi() *gin.Engine {
	// Create a router
	router := gin.Default()

	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1"}, // TODO: Change this to wherever the front-end is located!
		AllowMethods:     []string{"OPTIONS", "PUT", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	root := router.Group(globals.ApiEndpoint)

	// Add the proxy routes
	for _, router := range routers {
		router.AddToGroup(root)
	}

	statikFS, err := fs.New()
	if err != nil {
		logger.Log.Fatal().Err(err)
	}

	// Serve the Swagger UI files in the root of the api
	// TODO: [7/24/2022] Use Pakr or Statik to bundle non-golang files into the binary
	root.StaticFS("swagger", statikFS)

	return router
}

func ParseFlags() {
	// Set the logging level to debug
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Starts the API
	if *runApi {
		// Serve the API
		api := setupApi()
		api.Run(fmt.Sprintf("%s:%s", globals.ApiHost, globals.ApiPort))
	}
}
