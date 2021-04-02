/*
This is a template for all the service plugins that implements
all the requirements from the `Service` interface.

The name of the file must be called as the service plus a `d` i.e.
`templated.go`, otherwise it won't be discovered.

TIP: To create your own plugin...
	1. Create a folder with the name of your plugin in the `~/riotpot` folder.
	2. Create a `go` file with the same name.
	3. `ctrl + F` > replace > "Template" with "<your pluggin name>"

	4. Test it and place it in the `~/env/.env` file.
*/
package main

import (
	"riotpot/services"
)

var Name string

// Place here the name of the function which returns the service interface.
// This name will be used as a lookup symbol.
func init() {
	Name = "Templated"
}

// The function must be capitalize or exported, and return a `Service`
// interface compatible struct.
func Templated() services.Service {
	return &Template{
		id:   123,
		name: Name,
		port: 0,
		host: "localhost",
	}
}

// Template structure, defines the common fields the structure
// will use.
type Template struct {
	// it is recommended to include some kind of identity
	// for the service.
	id   int
	name string

	// declare here any other variable relevant for the
	// service to run. This are merely examples...
	port int
	host string
}

func (e *Template) Init(map[string]interface{}) {}

func (e *Template) Run() error {
	var err error
	// Place here your logic...
	return err
}

func (e *Template) Stop() error {
	var err error
	// Place here your logic...
	return err
}

func (e *Template) Restart() error {
	var err error
	// Place here your logic...
	return err
}

func (e *Template) Status() error {
	var err error
	// Place here your logic...
	return err
}

func (e *Template) Logger(ch chan<- error) (services.Logger, error) {
	var (
		logger services.Logger
		err    error
	)
	return logger, err
}
