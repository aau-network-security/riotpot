package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apiProxy "github.com/riotpot/api/proxy"
	apiService "github.com/riotpot/api/service"
	"github.com/riotpot/internal/proxy"
	"github.com/riotpot/pkg/services"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	// Create a router
	router := gin.Default()
	group := router.Group("/api/")
	// Add the proxy routes
	apiProxy.ProxiesRouter.AddToGroup(group)
	apiService.ServicesRouter.AddToGroup(group)

	return router
}

func TestApiProxy(t *testing.T) {

	// Add a proxy to the manager
	pe, err := proxy.Proxies.CreateProxy(proxy.TCP, 8080)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response with the proxy
	mockResponse := fmt.Sprintf(`[{"id":%s,"port":%d,"protocol":"%s","status":false,"service":null}]`, pe.GetID(), pe.GetPort(), pe.GetProtocol())

	router := SetupRouter()
	req, err := http.NewRequest("GET", "/api/proxies/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
}

func TestApiService(t *testing.T) {
	// Add a proxy to the manager
	se, err := services.Services.CreateService("Test Service", 8080, proxy.TCP, "localhost")
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response with the proxy
	mockResponse := fmt.Sprintf(`[{"id":"%s","name":"%s","port":%d,"host":"%s","protocol":"%s","locked":%t}]`, se.GetID(), se.GetName(), se.GetPort(), se.GetHost(), se.GetProtocol(), se.IsLocked())

	router := SetupRouter()
	req, err := http.NewRequest("GET", "/api/services/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
}
