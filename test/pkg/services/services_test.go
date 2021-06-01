package test_services

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"context"
	"log"
	"time"

	"github.com/riotpot/pkg/services"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
)

var ser services.Service

const (
	ECHOD = "Echod"
	COAP  = "Coapd"
)

func init() {
	servs := []string{
		"../../../pkg/plugin/coapd/plugin.so",
	}

	services := services.Services{}
	services.AutoRegister(servs)

	serviceToTest := COAP

	ser = services.Get(serviceToTest)
	if ser == nil {
		fmt.Print("Service not found")
	}

	fmt.Printf("[+] Ready to run %s\n", serviceToTest)

	// Initiate the service
	go ser.Run()
	fmt.Printf("[+] Service %s Running...\n", serviceToTest)
}

func TestStopServiceRunning(t *testing.T) {
	// Log the current status of the service
	status := ser.Status()
	t.Log(status)

	// checks if the service is running.
	// this is also done by the function `Stop` though
	if status == "Running" {
		ser.Stop()
		t.Log(ser.Status())
	} else {
		t.Error("Service is not running, restarting")
		ser.Restart()

		// Note! this might cause a race condition
		t.Log(ser.Status())
	}
}

func TestCoapDiscovery(t *testing.T) {
	co, err := udp.Dial("localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	path := "/ps"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Send GET Discovery /ps?ct=0
	co.Get(ctx, fmt.Sprintf("%s?ct=0", path))

	// Send GET
	co.Get(ctx, path)

	// Send POST "hello!"
	co.Post(ctx, path, message.TextPlain, bytes.NewReader([]byte("hello!")))
}

func TestCoapObserver(t *testing.T) {

	sync := make(chan bool)
	co, err := udp.Dial("localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	num := 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	obs, err := co.Observe(ctx, "/some/path", func(req *pool.Message) {
		log.Printf("Got %+v\n", req)
		num++
		if num >= 10 {
			sync <- true
		}
	})
	if err != nil {
		log.Fatalf("Unexpected error '%v'", err)
	}
	<-sync
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	obs.Cancel(ctx)
}

func TestNetServiceRunning(t *testing.T) {
	servers := []struct {
		protocol string
		port     string
	}{
		// List of servers to be tested
		{"tcp", ":7"},
	}

	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		// Mock data
		{"Sending a simple request returns result", []byte("hello world\n"), []byte("Request received: hello world")},
		{"Sending another simple request works", []byte("goodbye world\n"), []byte("Request received: goodbye world")},
	}

	// Iterate through the servers that we want to test
	for _, serv := range servers {

		// iterate the testing data
		for _, tc := range tt {

			// ping the server at the given port
			conn, err := net.Dial(serv.protocol, serv.port)
			if err != nil {
				t.Error("could not connect to server: ", err)
			}
			defer conn.Close()

			// Log the test
			t.Log(tc.test)

			// write some data to the server
			if _, err := conn.Write(tc.payload); err != nil {
				t.Errorf("Could not write payload to server: %v", err)
			}

			// read the response from the server and compare it to the desired
			// response.
			out := make([]byte, 1024)
			if _, err := conn.Read(out); err == nil {
				if bytes.Equal(out, tc.want) {
					t.Log("Messages were equal")
				} else {
					t.Errorf("Unexpected response: %s", out)
				}
			} else {
				t.Error("Could not read from connection")
			}
		}
	}
}
