package main

import (
	"fmt"
	"net"

	"github.com/riotpot/pkg/services"
	"github.com/riotpot/tools/errors"
	"github.com/xiegeo/modbusone"
)

var Name string

const size = 0x10000

var (
	discretes        [size]bool
	coils            [size]bool
	inputRegisters   [size]uint16
	holdingRegisters [size]uint16
)

func init() {
	Name = "Modbusd"
}

func Modbusd() services.Service {
	mx := services.MixinService{
		Name:     Name,
		Port:     502, // Also 1502 in some cases
		Running:  make(chan bool, 1),
		Protocol: "tcp",
	}

	return &Modbus{
		mx,
	}
}

type Modbus struct {
	services.MixinService
}

func (m *Modbus) Run() (err error) {
	// convert the port number to a string that we can use it in the server
	var port = fmt.Sprintf(":%d", m.Port)

	// start a service in the `echo` port
	listener, err := net.Listen("tcp", port)
	errors.Raise(err)

	// create the channel for stopping the service
	m.StopCh = make(chan int, 1)

	// serve the app
	server := modbusone.NewTCPServer(listener)
	m.serve(server)

	// update the status of the service
	m.Running <- true

	// block here until we receive an stopping signal in the channel
	<-m.StopCh

	// close the server
	defer server.Close()

	// Close the channel for stopping the service
	fmt.Print("[x] Service stopped...\n")
	close(m.StopCh)

	return
}

func (m *Modbus) serve(server modbusone.ServerCloser) {
	fmt.Printf("[%s] Started listenning for connections in port %d\n", Name, m.Port)

	// since we don't have to handle the connection as the server already
	// has a handler, we just run the server in another goroutine
	go func() {
		err := server.Serve(m.handler("server"))
		errors.Raise(err)
	}()
}

// Simple handler for the Modbus functions.
// We will send adequate responses for each of them, however, we will store
// any information comming to the honeypot in the process.
// Modbus works with six different functions with their own individual code:
// 	1: read coils
// 	2: read discrete inputs
// 	3: read holding registers
// 	4: read input registers
// 	5: force/write single coil
// 	6: preset/write single holding register
// 	15: miltiple 5
// 	16: multiple 6
// Generarly, we will either read a quantity of an unsigned int 16, or write
// booleans (0 or 1) plus the address.
func (m *Modbus) handler(name string) modbusone.ProtocolHandler {
	return &modbusone.SimpleHandler{

		// Discrete Inputs
		ReadDiscreteInputs: func(address, quantity uint16) ([]bool, error) {
			return discretes[address : address+quantity], nil
		},
		WriteDiscreteInputs: func(address uint16, values []bool) error {
			for i, v := range values {
				discretes[address+uint16(i)] = v
			}
			return nil
		},

		// Coils
		ReadCoils: func(address, quantity uint16) ([]bool, error) {
			return coils[address : address+quantity], nil
		},
		WriteCoils: func(address uint16, values []bool) error {
			for i, v := range values {
				coils[address+uint16(i)] = v
			}
			return nil
		},

		// Registers
		ReadInputRegisters: func(address, quantity uint16) ([]uint16, error) {
			return inputRegisters[address : address+quantity], nil
		},
		WriteInputRegisters: func(address uint16, values []uint16) error {
			for i, v := range values {
				inputRegisters[address+uint16(i)] = v
			}
			return nil
		},

		// Holding registers
		ReadHoldingRegisters: func(address, quantity uint16) ([]uint16, error) {
			return holdingRegisters[address : address+quantity], nil
		},
		WriteHoldingRegisters: func(address uint16, values []uint16) error {
			for i, v := range values {
				holdingRegisters[address+uint16(i)] = v
			}
			return nil
		},

		// It might be that an "slave" wants to report an error
		// we want to gather that information as well.
		OnErrorImp: func(req modbusone.PDU, errRep modbusone.PDU) {
			fmt.Printf("error received: %v from req: %v\n", errRep, req)
		},
	}
}
