package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	api "github.com/riotpot/api/proxy"
	"github.com/riotpot/internal/proxy"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	// Create a router
	router := gin.Default()
	group := router.Group("/api/")
	// Add the proxy routes
	api.ProxiesRouter.Group(group)

	return router
}

func TestApi(t *testing.T) {

	// Add a proxy to the manager
	pe, _ := proxy.Proxies.CreateProxy(proxy.TCP, 8080)

	// Mock the response with the proxy
	mockResponse := fmt.Sprintf(`[{"id":%s,"port":%d,"protocol":"%s","status":false,"service":null}]`, pe.GetID(), pe.GetPort(), pe.GetProtocol())

	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/api/proxies/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
}
