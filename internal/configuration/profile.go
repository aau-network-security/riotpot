/*
This configuration package implements the structures
and method destined to the management of the profile
used with the riotpot installation.
*/
package configuration

import (
	"os"

	"github.com/riotpot/internal/cli"
	"github.com/riotpot/internal/greeting"
	errors "github.com/riotpot/tools/errors"

	"gopkg.in/yaml.v3"
)

func NewProfile() (p Profile, err error) {
	p = Profile{}
	err = p.Load("configs/samples/profile.yml")
	return
}

// Implements the User profile used with riotpot.
type Profile struct {
	Configuration

	// `Hello, World!`
	Greet greeting.Greet `yaml:"greet"`
	// Terminal configuration
	Terminal cli.Terminal

	// indicates the preferred mode to run riotpot
	// NOTE: currently not in use
	mode Options
}

// Load the configuration on the child.
func (conf *Profile) Load(path string) (err error) {
	data, err := os.ReadFile(path)
	errors.Raise(err)

	err = yaml.Unmarshal(data, &conf)
	errors.Raise(err)

	return err
}

// Stores the configuration into the given path in `.yml` format.
func (conf *Profile) Save(path string) (err error) {
	// marshal the content of the configuration into a `.yaml` document
	d, err := yaml.Marshal(&conf)
	errors.Raise(err)

	// Save to file.
	// Mode 640: https://chmodcommand.com/chmod-640/
	// Note: this truncates the file if it already exists !!!
	err = os.WriteFile(path, d, 0640)
	errors.Raise(err)

	return err
}

// Implements a selectable option
type Option struct {
	// Option number
	n int
	// Value of the option.
	value struct{}
}

// Implements an Option wrapper
type Options struct {
	// List of options in the wrapper
	o []Option

	// The currently selected option
	selected Option
}

// Method to add as many options as needed to the wrapper
func (ops *Options) Add(options ...Option) {
	ops.o = append(ops.o, options...)
}
