package messagebroker

import (
	"fmt"
)

// MessageBroker struct/kinda class
type MessageBroker struct {
	authenticatedClients map[string]*Client
	receivedPacket       chan Packet
	topics               map[string]map[*Client]bool
}

func (broker *MessageBroker) receive(client *Client) {
	for {
		message := Message{}
		_, err := message.read(client.socket)
		if err != nil {
			broker.disconnect(client)
			return
		}
		broker.receivedPacket <- Packet{
			sender:  client,
			message: message,
		}
	}
}

func (broker *MessageBroker) send(client *Client) {
	for envelope := range client.packetToSend {
		_, err := envelope.message.write(client.socket)
		if err != nil {
			break
		}
	}
	fmt.Println("Closing")
	close(client.packetToSend)
	client.socket.Close()
	broker.disconnect(client)
}

func (broker *MessageBroker) dispatchMessage() {
	for envelope := range broker.receivedPacket {
		if broker.handleMagicTopic(&envelope) {
			continue
		}

		subscribers := broker.topics[envelope.message.topic]
		for subscriber := range subscribers {
			subscriber.packetToSend <- envelope
		}
	}
}

func (broker *MessageBroker) handleMagicTopic(envelope *Packet) bool {
	switch envelope.message.topic {
	case "__SUB__":
		broker.onSubscribe(envelope)
	case "__UNSUB__":
		broker.onUnsubscribe(envelope)
	case "__AUTH__":
		broker.onAuth(envelope)
	default:
		return false
	}

	return true
}

func (broker *MessageBroker) onSubscribe(envelope *Packet) {
	topic := envelope.message.content
	subscriber := envelope.sender
	if broker.topics[topic] == nil {
		broker.topics[topic] = make(map[*Client]bool)
	}
	broker.topics[topic][subscriber] = true
	subscriber.topics[topic] = true
}

func (broker *MessageBroker) onUnsubscribe(envelope *Packet) {
	if subscribers := broker.topics[envelope.message.content]; subscribers != nil {
		delete(subscribers, envelope.sender)
	}
	delete(envelope.sender.topics, envelope.message.content)
}

func (broker *MessageBroker) disconnect(client *Client) {
	for topic := range client.topics {
		delete(broker.topics[topic], client)
		delete(client.topics, topic)
	}
}

func (broker *MessageBroker) onAuth(packet *Packet) {
	packet.sender.id = packet.message.content
	broker.authenticatedClients[packet.sender.id] = packet.sender
}
