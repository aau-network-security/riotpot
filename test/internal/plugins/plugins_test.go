package pkg

import (
	"testing"

	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/plugins"
	"github.com/riotpot/internal/services"
	"github.com/stretchr/testify/assert"
)

var (
	pluginPath = "../../../bin/plugins/*.so"
)

func TestLoadPlugins(t *testing.T) {
	pgs, err := plugins.GetPluginServices(pluginPath)
	if err != nil {
		lr.Log.Fatal().Err(err).Msgf("One or more services could not be found")
	}

	assert.Equal(t, 1, len(pgs))

	plg := pgs[0]
	i, ok := plg.(services.PluginService)
	if !ok {
		lr.Log.Fatal().Err(err).Msgf("Service is not a plugin")
	}
	go i.Run()
}

func TestNewPrivateKey(t *testing.T) {
	key := plugins.NewPrivateKey(plugins.DefaultKey)
	pem := key.GetPEM()

	assert.Equal(t, 1, len(pem))
}
