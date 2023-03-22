// This package implements mqtt sessions and packet tearing methods
// for mqtt connections.
// The packet information can be found here:
// 	- https://www.hivemq.com/blog/mqtt-essentials-part-6-mqtt-quality-of-service-levels/
//	- https://openlabpro.com/guide/mqtt-packet-format/
//  - http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718086
//
// Inspired on the mqtt broker connection handling of:
// https://github.com/luanjunyi/gossipd/
package main

import (
	"bytes"
	"io"
	"net"

	"github.com/riotpot/internal/logger"
)

func NewSession(conn net.Conn) *Session {
	return &Session{
		remote:           conn.RemoteAddr().String(),
		topics_available: []string{},
		subscriptions:    []Topic{},
	}
}

// Implements the loading of an mqtt session between
// the server and the client.
// It will simply stores the history of the session so we
// can respond properly to the subscriptions and publishings.
type Session struct {
	remote           string
	subscriptions    []Topic
	topics_available []string
}

type Topic struct {
	QoS  uint8
	Name string
}

// Method implementing the full reading of an mqtt connection packet.
func (s *Session) Read(conn net.Conn) (packet *Packet) {
	// read the fixed header from the packet
	f_header := NewFixedHeader(&conn)
	if f_header == nil {
		return nil
	}

	l := f_header.RemainingLength
	// create a buffer of exactly the same length as the remaining length of the packet
	// to check if there is any mistmatch between both.
	// Note: if so, there would be a malformed packet, and we want to capture it though,
	// but currently we just drop the connection!
	b := make([]byte, l)
	n, _ := io.ReadFull(conn, b)
	if uint32(n) > l {
		panic("mismatch between remaining length of the packet and the rest of the packet")
	}

	// read the packet based on the fixed header
	packet = NewPacket(f_header)
	packet.Decode(&conn)

	if len(packet.Topics) > 0 {
		s.subscribe(packet.Topics, packet.Topics_qos)
	}

	return
}

// Add the topic subscripton to the list
func (s *Session) subscribe(topics []string, qos []uint8) {
	for i, t := range topics {
		topic := Topic{
			Name: t,
			QoS:  qos[i],
		}
		s.subscriptions = append(s.subscriptions, topic)
	}
}

func (s *Session) Answer(p Packet, conn *net.Conn) {
	// copy the fixed header from the packet
	msg, err := p.EncodeAnswer()
	if err != nil {
		logger.Log.Error().Err(err)
		return
	}

	(*conn).Write(msg)
}

/* helping functions */

func boolToByte(b bool) (by byte) {
	by = byte(0)
	if b {
		by = byte(1)
	}
	return
}

func getUint8(b []byte, p *int) uint8 {
	*p += 1
	return uint8(b[*p-1])
}

func setUint8(val uint8, buf *bytes.Buffer) {
	buf.WriteByte(byte(val))
}

func getUint16(b []byte, p *int) uint16 {
	*p += 2
	return uint16(b[*p-2])<<8 + uint16(b[*p-1])
}

func setUint16(val uint16, buf *bytes.Buffer) {
	buf.WriteByte(byte(val & 0xff00 >> 8))
	buf.WriteByte(byte(val & 0x00ff))
}

func getString(b []byte, p *int) string {
	// each of this strings will be 2 bytes.
	length := int(getUint16(b, p))
	*p += length
	return string(b[*p-length : *p])
}
