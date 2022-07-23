package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/api"
	"github.com/riotpot/internal/validators"
	"github.com/riotpot/pkg/services"
)

type GetService struct {
	ID       string `json:"id" binding:"required" gorm:"primary_key"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
	Locked   bool   `json:"locked"`
}

type CreateService struct {
	Name     string `json:"name" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
}

// Routes
var (
	// General routes for the services
	servicesRoutes = []api.Route{
		// GET and POST services
		api.NewRoute("", "GET", getServices),
		api.NewRoute("", "POST", createService),
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
			ID:       serv.GetID(),
			Port:     serv.GetPort(),
			Name:     serv.GetName(),
			Host:     serv.GetHost(),
			Protocol: serv.GetProtocol(),
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

	sv, err := services.Services.CreateService(input.Name, input.Port, input.Protocol, input.Host)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret := NewService(sv)
	ctx.JSON(http.StatusOK, ret)
}

func patchService(ctx *gin.Context) {
	var errors []error

	// Small function to check whether the name is already taken
	validateName := func(name string) (n string, err error) {
		for _, service := range services.Services.GetServices() {
			if name == service.GetName() {
				err = fmt.Errorf("name already in use")
				return
			}
		}
		n = name
		return
	}

	// Validate the post request to patch the proxy
	var input GetService
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the service to update
	id := ctx.Param("id")
	sv, err := services.Services.GetService(id)
	if err != nil {
		errors = append(errors, err)
	}

	// Validate the port
	validPort, err := validators.ValidatePort(input.Port)
	if err != nil {
		errors = append(errors, err)
	}

	// Validate the name
	validName, err := validateName(input.Name)
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
	sv.SetLocked(input.Locked)

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
	log.Fatalf("Not implemented")
}
