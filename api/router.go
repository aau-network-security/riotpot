package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Proxy struct
type Proxy struct {
	ID   int `json:"id" binding:"required"`
	Port int `json:"port" binding:"required"`
}

// This function creates a new router that listens for new connections in
// the designated port
func NewRouter(port int) {
	// New default Gin router
	router := gin.Default()

	// Populate the routes
	api := router.Group("/api")
	{
		// Proxies
		api.GET("/proxies", getProxies)
	}

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })
	router.Run(fmt.Sprintf(":%d", port))
}

func getProxies(ctx *gin.Context) {
	// List of proxies
	var proxies []Proxy

	// Send the response with the proxies
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, proxies)

}
