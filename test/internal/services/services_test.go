package test_services

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"context"
	"time"

	"github.com/riotpot/internal/services"
	"github.com/stretchr/testify/assert"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
)

var (
	PORT = 5683
	HOST = "localhost"
)

func init() {
	services.Services.Start()
}

func TestCoapDiscovery(t *testing.T) {
	address := fmt.Sprintf("%s:%d", HOST, PORT)
	co, err := udp.Dial(address)
	if err != nil {
		t.Error(err)
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
	address := fmt.Sprintf("%s:%d", HOST, PORT)
	co, err := udp.Dial(address)
	if err != nil {
		t.Fatalf("Error dialing: %v", err)
	}
	num := 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	obs, err := co.Observe(ctx, "/some/path", func(req *pool.Message) {
		t.Logf("Got %+v\n", req)
		num++
		if num >= 10 {
			sync <- true
		}
	})
	if err != nil {
		t.Fatalf("Unexpected error '%v'", err)
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
				return
			}
			defer conn.Close()

			// Log the test
			t.Log(tc.test)

			// write some data to the server
			if _, err := conn.Write(tc.payload); err != nil {
				t.Errorf("Could not write payload to server: %v", err)
				return
			}

			// read the response from the server and compare it to the desired
			// response.
			out := make([]byte, 1024)
			if _, err := conn.Read(out); err == nil {
				assert.Equal(t, out, tc.want)
			} else {
				t.Error("could not read from connection")
			}
		}
	}
}
