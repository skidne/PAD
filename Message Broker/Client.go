package messagebroker

import "net"

// Client struct
type Client struct {
	id           string
	packetToSend chan Packet
	topics       map[string]bool
	socket       net.Conn
}
