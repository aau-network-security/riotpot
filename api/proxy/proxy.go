package proxy

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/api"
	"github.com/riotpot/internal/proxy"
	"github.com/riotpot/pkg/services"
)

// Structures used to serialize data
type Proxy struct {
	ID       string   `json:"id" gorm:"primary_key"`
	Port     int      `json:"port" binding:"required"`
	Protocol string   `json:"protocol" binding:"required"`
	Status   bool     `json:"status"`
	Service  *Service `json:"service"`
}

type Service struct {
	ID   string `json:"id" binding:"required" gorm:"primary_key"`
	Port int    `json:"port" binding:"required"`
	Name string `json:"name" binding:"required"`
	Host string `json:"host" binding:"required"`
}

// Routes
var (
	// Proxy Routes
	proxyRoutes = []api.Route{

		// GET and POST proxies
		api.NewRoute("", "GET", getProxies),
		api.NewRoute("", "POST", postProxy),

		// CRUD operations for each proxy
		api.NewRoute("proxy/:id", "GET", getProxy),
		api.NewRoute("proxy/:id", "PATCH", patchProxy),
		api.NewRoute("proxy/:id", "DELETE", delProxy),
	}
)

// Routers
var (
	ProxyRouter = api.NewRouter("proxies/", proxyRoutes, nil)
)

func newService(serv services.Service) (sv *Service) {
	if serv != nil {
		sv = &Service{
			ID:   serv.GetID(),
			Port: serv.GetPort(),
			Name: serv.GetName(),
			Host: serv.GetName(),
		}
	}

	return
}

func newProxy(px proxy.Proxy) *Proxy {
	serv := newService(px.Service())

	return &Proxy{
		ID:       px.ID(),
		Port:     px.Port(),
		Protocol: px.Protocol(),
		Status:   px.Alive(),
		Service:  serv,
	}
}

// TODO [7/17/2022]: Add filters to this method
// GET proxies registered
// Contains a filter to get proxies by port
func getProxies(ctx *gin.Context) {
	casted := []Proxy{}

	// Iterate through the proxies registered
	for _, px := range proxy.Proxies.Proxies() {
		// Serialize the proxy
		pr := newProxy(px)
		// Append the proxy to the casted
		casted = append(casted, *pr)
	}

	// Set the header and transform the struct to JSON format
	ctx.JSON(http.StatusOK, casted)
}

// POST a proxy by port ":port"
func postProxy(ctx *gin.Context) {
	// Validate the post request to create a new proxy
	var input Proxy
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the port from the parameters
	port, err := strconv.Atoi(ctx.Param("port"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new proxy
	pe, err := proxy.Proxies.CreateProxy(ctx.Param("protocol"), port)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the new proxy and return it as a response
	pr := newProxy(pe)
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
	pr := newProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

// TODO [7/17/2022]: Missing PATCH method
// PATCH proxy by port ":port"
// Can update port, protocol, status and service
func patchProxy(ctx *gin.Context) {
	// Validate the post request to create a new proxy
	var input Proxy
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
