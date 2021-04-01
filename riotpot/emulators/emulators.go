// This package provides multiple interfaces to load the emulators, validate them before running them
// and watching over their status
package emulators

import (
	"plugin"
	"riotpot/utils/errors"
)

// This struct provides the interface for the emulator to be loaded.
type Emulator struct {
	Id   int
	Name string
}

// Wrapper for the individual emulators.
type Emulators struct {
	Id int

	// List of emulators registered the wrapper
	emulators []Emulator
}

// Method used to append a new emulator to the list of the wrapper
func (es *Emulators) Register(emulator Emulator) {
	es.emulators = append(es.emulators, emulator)
}

// This function utilizes a list of starting emulators
// to create and register new emulators.
//
//	Note: This function does not discern between new and already running emulators!
func (es *Emulators) AutoRegister(emulators []string) {
	// iterate through the slice of emulators
	for _, emu := range emulators {
		// create a new emulator and register it.
		emulator := Emulator{
			Id:   len(es.emulators),
			Name: emu,
		}
		es.Register(emulator)
	}
}

// Function to get an stored plugin.
// based on: https://echorand.me/posts/getting-started-with-golang-plugins/
func runPlugin(path string, reference string) {

	// Open the plugin within the path
	pg, err := plugin.Open(path)
	errors.Raise(err)

	// check if the reference exists in the plugin
	rf, err := pg.Lookup(reference)
	errors.Raise(err)

	rf.(func())()
}
