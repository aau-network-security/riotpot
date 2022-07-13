package api

import (
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
	group := router.Group("/api")
	// Add the proxy routes
	group = api.ProxyRouter.Group(group)

	return router
}

func TestApi(t *testing.T) {
	mockResponse := `[{"id":0,"port":8080,"protocol":"tcp","status":false,"service":null}]`
	// Add a proxy to the manager
	proxy.Proxies.CreateProxy(proxy.TCP, 8080)
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/api/proxies/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
}
