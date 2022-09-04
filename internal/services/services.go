// This package provides multiple interfaces to load the services, validate them before running them
// and watching over their status
package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/validators"
)

// Interface used by every service plugin that offers a service. At the very least, every plugin
// must contain the set of methods and attributes from this interface.
// It is up to the plugin to determine the implementation of these methods for the most part.
type Service interface {
	// Get attributes from the structure
	GetID() string
	GetName() string
	GetNetwork() globals.Network
	GetInteraction() globals.Interaction
	GetPort() int
	GetAddress() string
	GetHost() string
	IsLocked() bool

	// Setters
	SetPort(port int) (int, error)
	SetName(name string)
	SetHost(host string)
	SetLocked(locked bool) (bool, error)
}

// Implements a mixin service that can be used as a base for any other service `struct` type.
type AbstractService struct {
	// require the methods described by `Service` on loading
	Service

	id          uuid.UUID
	name        string
	network     globals.Network
	port        int
	host        string
	locked      bool
	interaction globals.Interaction
}

// Getters
func (as *AbstractService) GetID() string {
	return as.id.String()
}

func (as *AbstractService) GetName() string {
	return as.name
}

func (as *AbstractService) GetNetwork() globals.Network {
	return as.network
}

func (as *AbstractService) GetInteraction() globals.Interaction {
	return as.interaction
}

func (as *AbstractService) GetPort() int {
	return as.port
}

func (as *AbstractService) GetHost() string {
	return as.host
}

func (as *AbstractService) GetAddress() string {
	return fmt.Sprintf("%s:%d", as.host, as.port)
}

func (as *AbstractService) IsLocked() bool {
	return as.locked
}

// Setters
func (as *AbstractService) SetPort(port int) (p int, err error) {
	err = validators.ValidatePortNumber(port)
	if err != nil {
		return
	}

	p = port
	as.port = port
	return
}

func (as *AbstractService) SetName(name string) {
	as.name = name
}

func (as *AbstractService) SetHost(host string) {
	as.host = host
}

func (as *AbstractService) SetLocked(locked bool) (bool, error) {
	as.locked = locked
	return as.locked, nil
}

// Implementation of a plugin-based service
// These services are stored localy as binary files that are mounted into the
// application as symbols that can be called
type PluginService interface {
	Run() error
}

type PluginServiceItem struct {
	PluginService
	service *AbstractService
}

func (aps *PluginServiceItem) GetID() string {
	return aps.service.GetID()
}

func (aps *PluginServiceItem) GetAddress() string {
	return aps.service.GetAddress()
}

func (aps *PluginServiceItem) GetNetwork() globals.Network {
	return aps.service.GetNetwork()
}

func (aps *PluginServiceItem) GetInteraction() globals.Interaction {
	return aps.service.GetInteraction()
}

func (aps *PluginServiceItem) GetPort() int {
	return aps.service.GetPort()
}

func (aps *PluginServiceItem) GetName() string {
	return aps.service.GetName()
}

func (aps *PluginServiceItem) GetHost() string {
	return aps.service.GetHost()
}

func (aps *PluginServiceItem) IsLocked() bool {
	return true
}

func (aps *PluginServiceItem) SetPort(port int) (p int, err error) {
	return aps.service.SetPort(port)
}

func (aps *PluginServiceItem) SetName(name string) {
	aps.service.SetName(name)
}

func (aps *PluginServiceItem) SetHost(host string) {
	aps.service.SetHost(host)
}

func (aps *PluginServiceItem) SetLocked(locked bool) (bool, error) {
	return true, fmt.Errorf("the lock status of this service can not change")
}

func NewService(name string, port int, network globals.Network, host string, interaction globals.Interaction) *AbstractService {
	return &AbstractService{
		id:          uuid.New(),
		name:        name,
		port:        port,
		network:     network,
		host:        host,
		interaction: interaction,
	}
}

// Simple constructor for plugin services
func NewPluginService(name string, port int, network globals.Network) Service {
	return &PluginServiceItem{
		service: NewService(name, port, network, "localhost", globals.Low),
	}
}
