package globals

import (
	"strconv"

	"github.com/riotpot/tools/environ"
)

type (
	Status      int8
	Network     int8
	Interaction int8
)

// Proxy status
const (
	// Status representing a stopped proxy
	StoppedStatus Status = iota
	// Status representing a running proxy
	RunningStatus

	// Value used for stopped status
	StoppedStatusValue = "stopped"
	// Value used for running status
	RunningStatusValue = "running"
)

func (s Status) String() string {
	switch s {
	case StoppedStatus:
		return StoppedStatusValue
	case RunningStatus:
		return RunningStatusValue
	}

	return strconv.Itoa(int(s))
}

func ParseStatus(status string) (st Status, err error) {
	switch status {
	case StoppedStatus.String():
		return StoppedStatus, nil
	case RunningStatus.String():
		return RunningStatus, nil
	}

	// Attempt to convert the status into an integer
	i, err := strconv.Atoi(status)
	if err != nil {
		return
	}

	// Return the status by the number
	return Status(i), nil
}

// Network Protocol
const (
	// TCP network protocol
	TCP Network = iota
	// UDP network protocol
	UDP

	// Value for TCP
	TCPValue = "tcp"
	// Value for UDP
	UDPValue = "udp"
)

func (n Network) String() string {
	switch n {
	case TCP:
		return TCPValue
	case UDP:
		return UDPValue
	}

	return strconv.Itoa(int(n))
}

func ParseNetwork(network string) (nt Network, err error) {
	switch network {
	case TCP.String():
		return TCP, nil
	case UDP.String():
		return UDP, nil
	}

	// If this point is reached, attempt to convert the string to integer
	i, err := strconv.Atoi(network)
	if err != nil {
		return
	}

	// Return the network and the error
	return Network(i), nil
}

// Interaction Level
const (
	// Low level
	Low Interaction = iota
	High

	// Value for Low
	LowValue  = "low"
	HighValue = "high"
)

func (i Interaction) String() string {
	switch i {
	case High:
		return HighValue
	case Low:
		return LowValue
	}

	return strconv.Itoa(int(i))
}

func ParseInteraction(interaction string) (in Interaction, err error) {
	switch interaction {
	case High.String():
		return High, nil
	case Low.String():
		return Low, nil
	}

	i, err := strconv.Atoi(interaction)
	if err != nil {
		return
	}

	return Interaction(i), nil
}

// API
var (
	// Root of the API endpoint
	ApiEndpoint string = environ.Getenv("API_ROOT", "/api/")
	// Address of the API
	ApiHost string = environ.Getenv("API_HOST", "localhost")
	// Port in where the API is listening
	ApiPort string = environ.Getenv("API_PORT", "2022")
)

// Database
var (
	// Database username
	DbUsername string = environ.Getenv("DB_USER", "username")
	// Database user password
	DbPassword string = environ.Getenv("DB_PASS", "password")
	// Database host
	DbHost string = environ.Getenv("DB_HOST", "localhost")
	// Database port
	DbPort string = environ.Getenv("DB_PORT", "5432")
	// Database name of the targeted database
	DbName string = environ.Getenv("DB_Name", "db")
)
