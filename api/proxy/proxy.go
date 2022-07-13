package proxy

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/api"
	"github.com/riotpot/internal/proxy"
	"github.com/riotpot/pkg/services"
)

// Structures used to serialize data
type Proxy struct {
	ID       int      `json:"id" binding:"required" gorm:"primary_key"`
	Port     int      `json:"port" binding:"required"`
	Protocol string   `json:"protocol" binding:"required"`
	Status   bool     `json:"status"`
	Service  *Service `json:"service"`
}

type Service struct {
	ID   int    `json:"id" binding:"required" gorm:"primary_key"`
	Port int    `json:"port" binding:"required"`
	Name string `json:"name" binding:"required"`
	Host string `json:"host" binding:"required"`
}

// Routes
var (
	// Proxy Routes
	proxyRoutes = []api.Route{
		// Get proxies
		api.NewRoute("/", "GET", getProxies),
		// Get a proxy by port
		api.NewRoute("/proxy/:port", "GET", getProxy),
		// Post proxy by port
		api.NewRoute("/proxy/:port", "POST", postProxy),
		// Delete a proxy by port
		api.NewRoute("/proxy/:port", "DELETE", delProxy),
		// Patch (Not update) a proxy by port
		api.NewRoute("/proxy/:port", "PATCH", patchProxy),
	}
)

// Routers
var (
	ProxyRouter = api.NewRouter("/proxies", proxyRoutes, nil)
)

func newService(serv services.Service) (sv *Service) {
	if serv != nil {
		sv = &Service{
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
		Port:     px.Port(),
		Protocol: px.Protocol(),
		Status:   px.Alive(),
		Service:  serv,
	}
}

// GET proxies registered
// Contains a filter to get proxies by port
func getProxies(ctx *gin.Context) {
	casted := []Proxy{}

	// Iterate through the proxies registered
	for _, px := range proxy.Proxies.Proxies() {
		// Serialize the proxy
		pr := newProxy(px)
		// Append the proxy tot he casted
		casted = append(casted, *pr)
	}

	// Set the header and transform the struct to JSON format
	ctx.JSON(http.StatusOK, casted)

}

// GET proxy by port ":port"
func getProxy(ctx *gin.Context) {
	port, err := strconv.Atoi(ctx.Param("port"))
	if err != nil {
		log.Fatal(err)
	}

	// Get the proxy
	px, err := proxy.Proxies.GetProxy(port)

	// If the proxy could not be found, let the user know
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the proxy and send it as a response
	pr := newProxy(px)
	ctx.JSON(http.StatusOK, pr)
}

// POST a proxy by port ":port"
func postProxy(ctx *gin.Context) {}

// DELETE a proxy by port ":port"
func delProxy(ctx *gin.Context) {}

// PATCH proxy by port ":port"
// Can update port, protocol, status and service
func patchProxy(ctx *gin.Context) {}

// Routers
/*
var (
	ProxyRouter = &api.AbstractRouter{
		path: "proxy",
		routes: []Route{
			GetProxies,
			GetProxy,
		},
	}
)
*/

/*
func getProxies(ctx *gin.Context) {
	// List of proxies
	var proxies []Proxy

	// Send the response with the proxies
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, proxies)

}
*/
