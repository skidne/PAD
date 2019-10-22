package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Message struct {
	Command string
	Topic   string
	Content string
}

var clientCount = 0
var allClients = make(map[net.Conn]int)
var subsribedClients = make(map[string][]net.Conn)
var newConnections = make(chan net.Conn)

func connect() {
	server, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()
}

func subscribe(client net.Conn, topic string) {
	subsribedClients[topic] = append(subsribedClients[topic], client)
}

func filter(list []net.Conn, toRemove net.Conn) (ret []net.Conn) {
	for _, s := range list {
		if s != toRemove {
			ret = append(ret, s)
		}
	}
	return
}

func unsubscribe(client net.Conn, topic string) {
	subsribedClients[topic] = filter(subsribedClients[topic], client)
}

func publish(topic string, content string) {
	for t, conn := range subsribedClients {
		if t == topic {
			for i := 0; i < len(conn); i++ {
				go func(conn net.Conn, message string) {
					msg := &Message{Topic: topic, Content: message}
					json, _ := json.Marshal(msg)
					conn.Write(json)
				}(conn[i], content)
			}
		}
	}
}

func parseCommands(conn net.Conn, msg Message) {
	switch msg.Command {
	case "%sub%":
		subscribe(conn, msg.Topic)
	case "%unsub%":
		unsubscribe(conn, msg.Topic)
	case "%pub%":
		publish(msg.Topic, msg.Content)
	}
}

func read(conn net.Conn, clientId int) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		var msg Message
		json.Unmarshal([]byte(incoming), &msg)
		parseCommands(conn, msg)
		log.Printf("New message: %s", incoming)
		conn.Write([]byte("Message received by broker."))
	}
}

func main() {
	connect()

	for {
		select {
		case conn := <-newConnections:
			log.Printf("Serving new client, #%d", clientCount)

			allClients[conn] = clientCount
			clientCount++

			go read(conn, allClients[conn])
		}
	}
}
