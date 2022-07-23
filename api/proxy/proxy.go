package proxy

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/api"
	"github.com/riotpot/api/service"
	"github.com/riotpot/internal/proxy"
	"github.com/riotpot/internal/validators"
	"github.com/riotpot/pkg/services"
)

// Structures used to serialize data:
type GetProxy struct {
	ID       string              `json:"id" binding:"required" gorm:"primary_key"`
	Port     int                 `json:"port"`
	Protocol string              `json:"protocol"`
	Status   int                 `json:"status"`
	Service  *service.GetService `json:"service"`
}

type CreateProxy struct {
	Port     int    `json:"port" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
}

// Routes
var (

	// General routes for proxies
	proxiesRoutes = []api.Route{
		// GET and POST proxies
		api.NewRoute("", "GET", getProxies),
		api.NewRoute("", "POST", createProxy),
	}

	// Routes to manipulate a proxy
	proxyRoutes = []api.Route{
		// CRUD operations for each proxy
		api.NewRoute("", "GET", getProxy),
		api.NewRoute("", "PATCH", patchProxy),
		api.NewRoute("", "DELETE", delProxy),
		api.NewRoute("/status", "POST", changeProxyStatus),
	}
)

// Routers
var (
	// Proxies
	ProxiesRouter = api.NewRouter("proxies/", proxiesRoutes, []api.Router{ProxyRouter})
	ProxyRouter   = api.NewRouter(":id/", proxyRoutes, []api.Router{service.ServiceRouter})
)

func NewProxy(px proxy.Proxy) *GetProxy {
	serv := service.NewService(px.GetService())

	return &GetProxy{
		ID:       px.GetID(),
		Port:     px.GetPort(),
		Protocol: px.GetProtocol(),
		Status:   px.GetStatus(),
		Service:  serv,
	}
}

// TODO [7/17/2022]: Add filters to this method
// GET proxies registered
// Contains a filter to get proxies by port
func getProxies(ctx *gin.Context) {
	casted := []GetProxy{}

	// Iterate through the proxies registered
	for _, px := range proxy.Proxies.GetProxies() {
		// Serialize the proxy
		pr := NewProxy(px)
		// Append the proxy to the casted
		casted = append(casted, *pr)
	}

	// Set the header and transform the struct to JSON format
	ctx.JSON(http.StatusOK, casted)
}

// POST a proxy by port ":port"
func createProxy(ctx *gin.Context) {
	// Validate the post request to create a new proxy
	var input CreateProxy
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new proxy
	pe, err := proxy.Proxies.CreateProxy(input.Protocol, input.Port)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the new proxy and return it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

func getProxy(ctx *gin.Context) {
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the proxy and send it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

// Can update:
// port, status and service
func patchProxy(ctx *gin.Context) {
	var errors []error

	// Validate the post request to patch the proxy
	var input GetProxy
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the proxy to update
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		errors = append(errors, err)
	}

	// Validate the port and the service
	validServ, err := services.Services.GetService(input.Service.ID)
	if err != nil {
		errors = append(errors, err)
	}

	validPort, err := validators.ValidatePort(input.Port)
	if err != nil {
		errors = append(errors, err)
	}

	// If there are errors in the list, send a message to the client and return
	if len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Update the Port
	pe.SetPort(validPort)
	// Update the service
	pe.SetService(validServ)

	// Serialize the proxy and send it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

// DELETE registered proxy
func delProxy(ctx *gin.Context) {
	id := ctx.Param("id")

	err := proxy.Proxies.DeleteProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the proxy and send it as a response
	ctx.JSON(http.StatusOK, gin.H{"success": "Proxy deleted"})
}

// POST request to change the status of the proxy
func changeProxyStatus(ctx *gin.Context) {
	// Get the proxy using the id
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the status string to integer
	status, err := strconv.Atoi(ctx.Param("status"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attempt to change the status
	switch status {
	case proxy.ALIVE:
		err = pe.Start()
	case proxy.DEAD:
		err = pe.Stop()
	default:
		err = fmt.Errorf("status not allowed")
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the status and send it as the response
	ctx.JSON(http.StatusOK, gin.H{"status": pe.GetStatus()})
}
