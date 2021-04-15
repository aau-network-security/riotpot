package test_services

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/riotpot/pkg/services"
)

var ser services.Service

func init() {
	servs := []string{
		"./echod.so",
	}

	services := services.Services{}
	services.AutoRegister(servs)

	serviceToTest := "Echod"

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
