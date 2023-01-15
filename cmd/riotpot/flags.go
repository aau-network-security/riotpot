/*
This package determines the flags set in the environment
*/

package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"github.com/riotpot/api"
	"github.com/riotpot/api/proxy"
	"github.com/riotpot/api/service"
	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/plugins"
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
	debug        = flag.Bool("debug", true, "Set log level to debug")
	runApi       = flag.Bool("api", true, "Whether to start the API")
	loadPlugins  = flag.Bool("plugins", true, "Whether to load the low-interaction honeypot plugins")
	allowedHosts = flag.String("whitelist", "http://127.0.0.1,http://localhost:3000", "List of allowed hosts to contact the API")
)

func setupApi(allowedHosts []string) *gin.Engine {
	// Create a router
	router := gin.Default()

	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedHosts,
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
	root.StaticFS("swagger", statikFS)

	return router
}

func ParseFlags() {
	// Set the logging level to debug
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Load the plugins
	if *loadPlugins {
		plugins.LoadPlugins()
	}

	// Starts the API
	if *runApi {
		// Serve the API
		whitelist := strings.Split(*allowedHosts, ",")
		api := setupApi(whitelist)

		apiAddress := fmt.Sprintf("%s:%s", globals.ApiHost, globals.ApiPort)
		api.Run(apiAddress)
	}
}
