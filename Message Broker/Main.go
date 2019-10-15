package messagebroker

import (
	"fmt"
	"net"
)

func (broker *MessageBroker) start() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", ":2000")
	if err != nil {
		fmt.Println(err)
	}

	go broker.dispatchMessage()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		client := &Client{
			socket:       connection,
			id:           "",
			packetToSend: make(chan Packet, 10),
			topics:       make(map[string]bool),
		}
		go broker.receive(client)
		go broker.send(client)
	}
}

func main() {
	broker := MessageBroker{
		topics:               make(map[string]map[*Client]bool),
		authenticatedClients: make(map[string]*Client),
		receivedPacket:       make(chan Packet, 10),
	}

	broker.start()
}
