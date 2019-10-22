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
var deadConnections = make(chan net.Conn)
var messages = make(chan string)

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
					_, err := conn.Write([]byte(message))

					if err != nil {
						deadConnections <- conn
					}
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

func main() {
	connect()

	for {
		select {
		case conn := <-newConnections:

			log.Printf("Serving new client, #%d", clientCount)

			allClients[conn] = clientCount
			clientCount++

			go func(conn net.Conn, clientId int) {
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
				}

				deadConnections <- conn

			}(conn, allClients[conn])

		case conn := <-deadConnections:
			log.Printf("Client %d disconnected", allClients[conn])
			delete(allClients, conn)
		}
	}
}
