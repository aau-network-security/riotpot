package main

import (
	"testing"
)

func TestEcho(t *testing.T) {
	// testing of the echo service
	service := Echod()

	// send the message to the server by typing in another terminal
	// $ nc localhost 7
	t.Log("For this test you will need a packet capture tool !!")
	t.Log("Please run the following command in another terminal:")
	t.Log("$ nc localhost 7")

	// run the Echo host in a blocking thread until the test is interrupted
	service.Run()
}
