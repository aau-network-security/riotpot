package main

import (
	"fmt"
	"io"
	"net"

	"github.com/riotpot/internal/logger"
	"github.com/riotpot/tools/arrays"
)

func NewFixedHeader(conn *net.Conn) *FixedHeader {
	b := make([]byte, 2)
	// check if there is a fixed header at all

	n, _ := io.ReadFull(*conn, b)
	if n != len(b) {
		logger.Log.Error().Msg("read header failed\n")
		return nil
	}

	// control header: Command type and Control flag
	ctr_header := b[0]
	rest := b[1]

	header := FixedHeader{
		// binary shifts after the control header
		MessageType:     uint8(ctr_header & 0xF0 >> 4),
		Duplicate:       ctr_header&0x08 > 0,
		QoS:             uint8(ctr_header & 0x06 >> 1),
		Retain:          ctr_header&0x01 > 0,
		RemainingLength: decodeVarLength(rest, conn),
	}

	return &header
}

type FixedHeader struct {
	// there are threee levels of QoS on MQTT protocol:
	// At most once (0) - Fire and forget
	// At least once (1)
	// Exactly once (2).
	QoS uint8

	// Type of the message, see `valType` func.
	MessageType uint8

	// Retain flag of the mqtt packet.
	// The broker stores the last retained message and the corresponding QoS for that topic
	Retain bool

	// used when re-publishing a message with QoS or 1 or 2
	Duplicate bool

	// The remaining length is the number of bytes following the length field,
	// includes variable length header length and payload lengths combined.
	RemainingLength uint32
}

func (f *FixedHeader) types() []string {
	return []string{
		"Reserved",    // f
		"CONNECT",     // client to server
		"CONNACK",     // server to client
		"PUBLISH",     // <->
		"PUBACK",      // <->
		"PUBREC",      // <->
		"PUBREL",      // <->
		"PUBCOMP",     // <->
		"SUBSCRIBE",   // client to server
		"SUBACK",      // server to client
		"UNSUBSCRIBE", // client to server
		"UNSUBACK",    // server to client
		"PINGREQ",     // client to server
		"PINGRESP",    // server to client
		"DISCONNECT",  // client to server
		"Reserved",    // f
	}
}

func (f *FixedHeader) TypeStr() string {
	return f.types()[f.MessageType]
}

func (f *FixedHeader) valToType(msgIndex uint8) (string, error) {
	types := f.types()

	if msgIndex > uint8(len(types)) {
		return "", fmt.Errorf("unknown type")
	}

	// return the string value of the message type index
	return types[msgIndex], nil
}

// Get the type of the response. This is only necessary for qos 1 or 2. The
// qos 0 does not need any reply
func (f *FixedHeader) responseType(msgIndex uint8, qos uint8) (t string, rspType uint8, err error) {
	msg, _ := f.valToType(f.MessageType)

	// we only respond to some messages, therefore we check if the message from the sender
	// is included in the `allowed` slice.
	allowed := []string{"CONNECT", "PUBLISH", "PUBREL", "SUBSCRIBE", "UNSUBSCRIBE", "PINGREQ"}

	if !arrays.Contains(allowed, msg) {
		err = fmt.Errorf("message type not allowed")
		return
	}

	switch msg {
	case "PUBLISH":
		if f.QoS == 1 {
			// for qos 1, we only send the publish acknowledgement
			t = "PUBACK"
			rspType = msgIndex + 1

		} else {
			// for qos 2, the process for responding is longer and includes a 3 messages
			// long conversation.
			t = "PUBREC"
			rspType = msgIndex + 2
		}
	default:
		// otherwise, the next message in the line
		// however, we need to be sure that the message is an ACK!
		rspType = msgIndex + 1
		t, err = f.valToType(rspType)
	}

	return
}

func decodeVarLength(cur byte, conn *net.Conn) uint32 {
	length := uint32(0)
	multi := uint32(1)

	for {
		// determine if we should continue reading
		// the section of the remaining length is divided into up to 4 bytes
		// in which each byte uses up to 7 bits for the length and the last
		// as a continuation bit, either 1 or 0. If set to 1, the next
		// byte will be a part of the length, otherwise it is the end.
		length += multi * uint32(cur&0x7f)
		// 	if there is no difference between the current
		// 	byte and 128, break the loop.
		if cur&0x80 == 0 {
			break
		}

		buf := make([]byte, 1)
		n, _ := io.ReadFull(*conn, buf)
		if n != 1 {
			panic("failed to read variable length in MQTT header")
		}
		cur = buf[0]
		multi *= 128
	}

	return length
}
