package main

import (
	"fmt"
	"io"
	"net"

	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
	"github.com/xiegeo/modbusone"
)

var Plugin string

const (
	name    = "Modbus"
	network = globals.TCP
	port    = 502
	size    = 0x10000
)

var (
	discretes        [size]bool
	coils            [size]bool
	inputRegisters   [size]uint16
	holdingRegisters [size]uint16
)

func init() {
	Plugin = "Modbusd"
}

func Modbusd() services.Service {
	mx := services.NewPluginService(name, port, network)

	handler := handler()

	return &Modbus{
		mx,
		handler,
	}
}

type Modbus struct {
	services.Service
	handler modbusone.ProtocolHandler
}

func (m *Modbus) Run() (err error) {

	// start a service in the `echo` port
	listener, err := net.Listen(m.GetNetwork().String(), m.GetAddress())
	logger.Log.Error().Err(err)

	// build a channel stack to receive connections to the service
	conn := make(chan net.Conn)
	m.serve(conn, listener)

	return
}

// Open the service and listen for connections
// inspired on https://gist.github.com/paulsmith/775764#file-echo-go
func (m *Modbus) serve(ch chan net.Conn, listener net.Listener) {
	// open an infinite loop to receive connections
	for {
		// Accept the client connection
		client, err := listener.Accept()
		if err != nil {
			return
		}
		defer client.Close()

		// push the client connection to the channel
		go m.handleSession(client)
	}
}

// Handle a connection made to the service
// We are rewriting this function so we can interact with the connection loop
// and rewrite some of the underlying function. Unfortunately, the loop
// is deep and contains a goroutine in it as main handler of the session connections.
func (m *Modbus) handleSession(conn net.Conn) {
	for {
		defer conn.Close()

		wec := func(conn net.Conn, bs []byte, req modbusone.PDU, err error) {
			writeTCP(conn, bs, modbusone.ExceptionReplyPacket(req, modbusone.ToExceptionCode(err)))
		}

		var rb []byte
		if modbusone.OverSizeSupport {
			rb = make(
				[]byte,
				modbusone.MBAPHeaderLength+modbusone.OverSizeMaxRTU+modbusone.TCPHeaderLength,
			)
		} else {
			rb = make([]byte, modbusone.MBAPHeaderLength+modbusone.MaxPDUSize)
		}

		for {
			n, err := readTCP(conn, rb)
			if err != nil {
				return
			}

			// load the payload from the packet i.e. everything after the header
			p := modbusone.PDU(rb[modbusone.MBAPHeaderLength:n])

			// validate the request, checks for errors on the code and the
			// length of the payload
			err = p.ValidateRequest()
			if err != nil {
				return
			}

			fc := p.GetFunctionCode()

			// initialize the data. It will be filled with the information
			// in the payload.
			var data []byte

			// Only two things can happen,
			// either read from the server or write to it.
			if fc.IsReadToServer() {
				data, err = m.handler.OnRead(p)
				if err != nil {
					wec(conn, rb, p, err)
					continue
				}
				writeTCP(conn, rb, p.MakeReadReply(data))
			} else if fc.IsWriteToServer() {
				data, err = p.GetRequestValues()
				if err != nil {
					wec(conn, rb, p, err)
					continue
				}
				err = m.handler.OnWrite(p, data)
				if err != nil {
					wec(conn, rb, p, err)
					continue
				}
				writeTCP(conn, rb, p.MakeWriteReply())
			}

		}
	}
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
func handler() modbusone.ProtocolHandler {
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
			logger.Log.Error().Msgf("error received: %v from req: %v\n", errRep, req)
		},
	}
}

// Copy pasted from https://github.com/xiegeo/modbusone/blob/797d647e237d97ab9d2bdad49bf42591ea7076f2/tcp_server.go#L19
// both functions are just packet handlers that read the content and headers.
// It is unnecessary to reimplement something already done, however, this function
// are now exported but still necessary for our connection handler.
func readTCP(r io.Reader, bs []byte) (n int, err error) {
	n, err = io.ReadFull(r, bs[:modbusone.TCPHeaderLength])
	if err != nil {
		return n, err
	}
	if bs[2] != 0 || bs[3] != 0 {
		return n, fmt.Errorf("MBAP protocol of %X %X is unknown", bs[2], bs[3])
	}
	l := int(bs[4])*256 + int(bs[5])
	if l <= 2 {
		return n, fmt.Errorf("MBAP data length of %v is too short, bs:%x", l, bs[:n])
	}
	if len(bs) < l+modbusone.TCPHeaderLength {
		return n, fmt.Errorf("MBAP data length of %v is too long", l)
	}
	n, err = io.ReadFull(r, bs[modbusone.TCPHeaderLength:l+modbusone.TCPHeaderLength])
	return n + modbusone.TCPHeaderLength, err
}

//writeTCP writes a PDU packet on TCP reusing the headers and buffer space in bs
func writeTCP(w io.Writer, bs []byte, pdu modbusone.PDU) (int, error) {
	l := len(pdu) + 1 //pdu + byte of slaveID
	bs[4] = byte(l / 256)
	bs[5] = byte(l)
	copy(bs[modbusone.MBAPHeaderLength:], pdu)
	return w.Write(bs[:len(pdu)+modbusone.MBAPHeaderLength])
}
