package main

import (
	"bytes"
	"net"
)

func NewPacket(fx *FixedHeader) (p *Packet) {
	return &Packet{
		FixedHeader: fx,
	}
}

type Packet struct {
	FixedHeader                                                                   *FixedHeader
	ProtocolName, TopicName, ClientId, WillTopic, WillMessage, Username, Password string
	ProtocolVersion                                                               uint8
	ConnectFlags                                                                  *ConnectionFlags
	KeepAliveTimer, MessageId                                                     uint16
	Data                                                                          []byte
	Topics                                                                        []string
	Topics_qos                                                                    []uint8
	ReturnCode                                                                    uint8
}

// Decode or Unmarshall the connection packet
func (p *Packet) Decode(conn *net.Conn) {
	buf := make([]byte, p.FixedHeader.RemainingLength)

	// get the type of the message as a string
	t := p.FixedHeader.TypeStr()

	// stablise where in the packet we are
	// this variable will be updated with each assign
	idx := 0

	switch t {
	case "CONNECT":
		// Protocol Name and version
		p.ProtocolName = getString(buf, &idx)
		p.ProtocolVersion = getUint8(buf, &idx)
		// Add the connection flags Now
		p.ConnectFlags = getConnectFlags(buf, &idx)
		p.KeepAliveTimer = getUint16(buf, &idx)
		p.ClientId = getString(buf, &idx)

		if p.ConnectFlags.WillFlag {
			p.WillTopic = getString(buf, &idx)
			p.WillMessage = getString(buf, &idx)
		}

		if p.ConnectFlags.Username && idx < len(buf) {
			p.Username = getString(buf, &idx)
		}

		if p.ConnectFlags.Password && idx < len(buf) {
			p.Password = getString(buf, &idx)
		}
	case "CONNACK":
		// the values of the returning code will be between 0 and 5, however,
		// this is a honeypot, therefore we want to log anything and everything!
		// Normally we would drop the connection because of a malformed message.
		// Anyway, the honeypot is not capable of stablishing connections by itself.
		idx += 1
		p.ReturnCode = uint8(getUint8(buf, &idx))
	case "PUBLISH":
		p.TopicName = getString(buf, &idx)

		// only the last 2 qos are restrictive in the amount
		// of messages that will be received and sent,
		// therefore the message will include an id for tracking.
		// same applies to "SUBSCRIBE" types.
		if p.FixedHeader.QoS > 0 {
			p.MessageId = getUint16(buf, &idx)
		}
		p.Data = buf[idx:]
		idx = len(buf)
	case "PUBACK", "PUBREC", "PUBCOMP", "UNSUBACK":
		// ACK, we only need to know the message ID for which we got the ack.
		p.MessageId = getUint16(buf, &idx)
	case "SUBSCRIBE":
		if p.FixedHeader.QoS > 0 {
			p.MessageId = getUint16(buf, &idx)
		}
		topics := make([]string, 0)
		topics_qos := make([]uint8, 0)

		// Each topic comes with a QoS and needs to be treated differently,
		// we only care to divide the packet and understand it to respond to it.
		for idx < len(buf) {
			topics = append(topics, getString(buf, &idx))
			topics_qos = append(topics_qos, getUint8(buf, &idx))
		}
		p.Topics = topics
		p.Topics_qos = topics_qos
	case "SUBACK":
		p.MessageId = getUint16(buf, &idx)
		topics_qos := make([]uint8, 0)
		for idx < len(buf) {
			topics_qos = append(topics_qos, getUint8(buf, &idx))
		}
		p.Topics_qos = topics_qos
	case "UNSUBSCRIBE":
		if p.FixedHeader.QoS > 0 {
			p.MessageId = getUint16(buf, &idx)
		}
		topics := make([]string, 0)
		for idx < len(buf) {
			topics = append(topics, getString(buf, &idx))
		}
		p.Topics = topics
	}

}

func (p *Packet) EncodeAnswer() (rsp []byte, err error) {

	var headerBuff, bodyBuff bytes.Buffer

	// transform the message into a string, just for better
	// readability.
	rspTypeStr, rspType, err := p.FixedHeader.responseType(
		p.FixedHeader.MessageType,
		p.FixedHeader.QoS,
	)

	// update the response type
	p.FixedHeader.MessageType = rspType

	// load the header in a buffer from where we can get the bytes
	p.setHeader(p.FixedHeader, &headerBuff)

	// this are the only response messages that we want to write any message.
	// Note: `PINGRESP` is missing in the list, but we still write the headers to
	// the client as a response.
	switch rspTypeStr {
	case "CONNACK":
		// The packet only contains a 0x0<Return Code Response from 0-5>
		// the only value that we consider is 0, in which we accept the connection.
		// In general, we will accept whatever, even though we only support 3.1
		bodyBuff.WriteByte(byte(0))
		setUint8(uint8(0), &bodyBuff)
	case "SUBACK":
		// Include the message ID to which we are responding.
		// Furthermore, we include for each of the topics
		// that the client wants to subscribe the desired QoS level
		// of the topic as a payload.
		setUint16(p.MessageId, &bodyBuff)
		for i := 0; i < len(p.Topics_qos); i += 1 {
			setUint8(p.Topics_qos[i], &bodyBuff)
		}

	case "UNSUBACK", "PUBCOMP", "PUBACK", "PUBREC":
		// Include the message id, thats the only info needed.
		setUint16(p.MessageId, &bodyBuff)
	}

	// encode the length of the payload
	p.encodeLength(uint32(bodyBuff.Len()), &headerBuff)
	//write the body next to the header
	headerBuff.Write(bodyBuff.Bytes())
	rsp = headerBuff.Bytes()
	return
}

func (p *Packet) encodeLength(length uint32, buf *bytes.Buffer) {
	if length == 0 {
		buf.WriteByte(byte(0))
		return
	}
	var lbuf bytes.Buffer
	for length > 0 {
		digit := length % 128
		length = length / 128
		if length > 0 {
			digit = digit | 0x80
		}
		lbuf.WriteByte(byte(digit))
	}
	blen := lbuf.Bytes()
	for i := 1; i <= len(blen); i += 1 {
		buf.WriteByte(blen[len(blen)-i])
	}
}

func (p *Packet) setHeader(header *FixedHeader, buf *bytes.Buffer) {
	val := byte(uint8(header.MessageType)) << 4
	val |= (boolToByte(header.Duplicate) << 3)
	val |= byte(header.QoS) << 1
	val |= boolToByte(header.Retain)
	buf.WriteByte(val)
}

func getConnectFlags(b []byte, p *int) *ConnectionFlags {
	bit := b[*p]
	*p += 1
	flags := ConnectionFlags{
		Username:     bit&0x80 > 0,
		Password:     bit&0x40 > 0,
		WillRetain:   bit&0x20 > 0,
		WillQoS:      uint8(bit & 0x18 >> 3),
		WillFlag:     bit&0x04 > 0,
		CleanSession: bit&0x02 > 0,
	}

	return &flags
}
