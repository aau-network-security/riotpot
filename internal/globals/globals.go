package globals

import (
	"strconv"
)

type (
	Status  int8
	Network int8
)

// Status
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

	i, err := strconv.Atoi(status)
	if err != nil {
		return
	}

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
