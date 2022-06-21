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
	"github.com/riotpot/pkg/services"
)

// Place here the name of the function which returns the service interface.
// This name will be used as a lookup symbol.
var Name string

func init() {
	Name = "Templated"
}

// The function must be capitalize or exported, and return a `Service`
// interface compatible struct.
func Templated() services.Service {
	mx := services.MixinService{
		Id:      123,
		Name:    Name,
		Port:    0,
		Host:    "localhost",
		Running: make(chan bool, 1),
	}

	return &Template{
		mx,
	}
}

// Template structure, implements the mixin containing common
// variables.
type Template struct {
	services.MixinService
}

func (e *Template) Run() error {
	var err error
	// Place here your logic...
	return err
}
