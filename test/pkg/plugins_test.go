package pkg

import (
	"testing"

	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
	"github.com/riotpot/pkg"
	"github.com/stretchr/testify/assert"
)

var (
	pluginPath = "../../bin/plugins/*.so"
)

func TestLoadPlugins(t *testing.T) {
	plugins, err := pkg.GetPluginServices(pluginPath)
	if err != nil {
		lr.Log.Fatal().Err(err).Msgf("One or more services could not be found")
	}

	plg := plugins[0]
	i, ok := plg.(services.PluginService)
	if !ok {
		lr.Log.Fatal().Err(err).Msgf("Service is not a plugin")
	}
	go i.Run()

	assert.Equal(t, 1, len(plugins))
}
