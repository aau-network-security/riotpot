/*
This is a template for all the service plugins that implements
all the requirements from the `Service` interface.

The name of the file must be called as the service plus a `d` i.e.
`templated.go`, otherwise it won't be discovered.

TIP: To create your own plugin...
	1. Create a folder with the name of your plugin in the `~/pkg/plugin` folder.
	2. Create a `go` file with the same name.
	3. `ctrl + F` > replace > "Template" with "<your pluggin name>"

	4. Test it and place it in the `~/env/.env` file.
*/
package main

import (
	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/services"
)

// Place here the name of the function which returns the service interface.
// This name will be used as a lookup symbol.
// NOTE: The name must be unique!
var Plugin string

var (
	name    string
	port    int
	network globals.Network
)

func init() {
	Plugin = "Templated"
}

// The function must be capitalize or exported, and return a `Service`
// interface compatible struct.
func Templated() services.PluginService {
	mx := services.NewPluginService(name, port, network)

	return &Template{
		mx,
	}
}

// Template structure, implements the mixin containing common
// variables.
type Template struct {
	services.Service
}

func (e *Template) Run() (err error) {
	// Place the plugin logic here
	return
}
