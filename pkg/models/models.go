// This package implements a series of useful model schemas to
// register in the database and use as template to create new
// entries in the database.
package models

import "time"

// Schema for a typical connection
type Connection struct {
	LocalAddress  string
	LocalPort     string
	RemoteAddress string
	RemotePort    string
	// payload sent/received
	Payload string
	// IP protocol
	Protocol string
	// the service running on the port
	Service string
	// wether the connection is from or to the server
	Incoming  bool
	Timestamp time.Time
}

func NewConnection() Connection {
	return Connection{
		// prepare the timestamp
		Timestamp: time.Now(),
	}
}
