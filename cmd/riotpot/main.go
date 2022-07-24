// Main Application
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/api/proxy"
	"github.com/riotpot/api/service"
)

var (
	ApiHost = "localhost"
	ApiPort = 2022
)

func API() *gin.Engine {
	// Create a router
	router := gin.Default()
	root := router.Group("/api/")

	// Add the proxy routes
	proxy.ProxiesRouter.AddToGroup(root)
	service.ServicesRouter.AddToGroup(root)

	// Serve the Swagger UI files in the root of the api
	// TODO: [7/24/2022] Use Pakr or Statik to bundle non-golang files into the binary
	root.Static("swagger", "api/swagger")

	return router
}

// `main` starts all the submodules containing the emulator services.
// It is the first function called when the application is run.
// It also acts as an orchestrator, which dictates the functioning of the application.
func main() {
	// Say Hi, don't be rude!
	fmt.Println("░▒▓███ RIoIPot ███▓▒░")

	// Serve the API
	api := API()
	api.Run(fmt.Sprintf("%s:%d", ApiHost, ApiPort))
}
