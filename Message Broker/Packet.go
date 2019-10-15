package messagebroker

// Packet struct
type Packet struct {
	sender  *Client
	message Message
}
