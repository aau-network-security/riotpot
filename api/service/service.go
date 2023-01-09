package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/api"
	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/proxy"
	"github.com/riotpot/internal/services"
	"github.com/riotpot/internal/validators"
)

type GetService struct {
	ID          string `json:"id" binding:"required" gorm:"primary_key"`
	Name        string `json:"name"`
	Port        int    `json:"port"`
	Host        string `json:"host"`
	Network     string `json:"network"`
	Locked      bool   `json:"locked"`
	Interaction string `json:"interaction"`
}

type CreateService struct {
	Name        string `json:"name" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Network     string `json:"network" binding:"required"`
	Interaction string `json:"interaction" binding:"required"`
}

type PatchService struct {
	Name string `json:"name" binding:"required"`
	Port int    `json:"port" binding:"required"`
	Host string `json:"host" binding:"required"`
}

type ServiceProxy struct {
	ID      string      `json:"id" binding:"required" gorm:"primary_key"`
	Port    int         `json:"port"`
	Network string      `json:"network"`
	Status  string      `json:"status"`
	Service *GetService `json:"service"`
}

// Routes
var (
	// General routes for the services
	servicesRoutes = []api.Route{
		// GET and POST services
		api.NewRoute("", "GET", getServices),
		api.NewRoute("", "POST", createService),
		api.NewRoute("new/", "POST", newServiceAndProxy),
	}

	// Routes to manipulate a service
	serviceRoutes = []api.Route{
		// CRUD operations for each service
		api.NewRoute("", "GET", getService),
		api.NewRoute("", "PATCH", patchService),
		api.NewRoute("", "DELETE", delService),

		// Get information about all the proxies this service is handling
		//api.NewRoute("proxies/", "GET", getServiceProxies),
	}
)

// Routers
var (
	// Services
	ServicesRouter = api.NewRouter("services/", servicesRoutes, []api.Router{ServiceRouter})
	ServiceRouter  = api.NewRouter(":id/", serviceRoutes, nil)
)

func NewService(serv services.Service) (sv *GetService) {
	if serv != nil {
		sv = &GetService{
			ID:          serv.GetID(),
			Port:        serv.GetPort(),
			Name:        serv.GetName(),
			Host:        serv.GetHost(),
			Network:     serv.GetNetwork().String(),
			Interaction: serv.GetInteraction().String(),
		}
	}
	return
}

func getServices(ctx *gin.Context) {
	casted := []GetService{}

	// Iterate through the services registered
	for _, sv := range services.Services.GetServices() {
		// Serialize the service
		ret := NewService(sv)
		// Append the service to the casted
		casted = append(casted, *ret)
	}

	// Set the header and transform the struct to JSON format
	ctx.JSON(http.StatusOK, casted)
}

func getService(ctx *gin.Context) {
	id := ctx.Param("id")
	sv, err := services.Services.GetService(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the service and send it as a response
	ret := NewService(sv)
	ctx.JSON(http.StatusOK, ret)
}

func createService(ctx *gin.Context) {
	// Validate the post request to patch the proxy
	var input CreateService
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nt, err := globals.ParseNetwork(input.Network)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	i, err := globals.ParseInteraction(input.Interaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sv, err := services.Services.CreateService(input.Name, input.Port, nt, input.Host, i)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret := NewService(sv)
	ctx.JSON(http.StatusOK, ret)
}

func newServiceAndProxy(ctx *gin.Context) {
	// Validate the post request to patch the proxy
	var input CreateService
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nt, err := globals.ParseNetwork(input.Network)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	i, err := globals.ParseInteraction(input.Interaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sv, err := services.Services.CreateService(input.Name, input.Port, nt, input.Host, i)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new proxy using the parameters from the service
	pe, err := proxy.Proxies.CreateProxy(nt, input.Port)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pe.SetService(sv)

	ret := &ServiceProxy{
		ID:      pe.GetID(),
		Port:    pe.GetPort(),
		Network: pe.GetNetwork().String(),
		Status:  pe.GetStatus().String(),
		Service: NewService(sv),
	}

	ctx.JSON(http.StatusOK, ret)
}

func patchService(ctx *gin.Context) {
	var errors []error

	// Small function to check whether the name is already taken
	validateName := func(id string, name string) (n string, err error) {
		for _, service := range services.Services.GetServices() {
			if name == service.GetName() && id != service.GetID() {
				err = fmt.Errorf("name already in use")
				return
			}
		}
		n = name
		return
	}

	// Validate the post request to patch the proxy
	var input PatchService
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the service to update
	id := ctx.Param("id")
	sv, err := services.Services.GetService(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the port
	validPort, err := validators.ValidatePort(input.Port)
	if err != nil {
		errors = append(errors, err)
	}

	// Validate the name
	validName, err := validateName(id, input.Name)
	if err != nil {
		errors = append(errors, err)
	}

	if ctx.Param("locked") != "" && !services.RemovableService(sv) {
		errors = append(errors, fmt.Errorf("the lock status of this service can not change"))
	}

	// If there are errors in the list, send a message to the client and return
	if len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Set the values
	sv.SetPort(validPort)
	sv.SetName(validName)
	sv.SetHost(input.Host)
	//sv.SetLocked(input.Locked)

	// Serialize the service and send it as a response
	ret := NewService(sv)
	ctx.JSON(http.StatusOK, ret)
}

func delService(ctx *gin.Context) {
	id := ctx.Param("id")

	err := services.Services.DeleteService(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the service and send it as a response
	ctx.JSON(http.StatusOK, gin.H{"success": "Proxy deleted"})
}

func getServiceProxies(ctx *gin.Context) {
	lr.Log.Fatal().Msg("Not implemented")
}
