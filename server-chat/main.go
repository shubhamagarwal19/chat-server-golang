package main

import (
	"fmt"
	"chat-server-golang/common"
	"io"
	"log"
	"net"
	"os"
	"bufio"
	"strings"
)

type Client struct {
	name string
	conn net.Conn
	chatWith string
}

var (
	connections []Client
)

func getInfo(bufc *bufio.Reader) []string {
	var userString []byte
	for len(userString) == 0 {
		userString, _, _ = bufc.ReadLine()
	}
	values := strings.Split(string(userString), ":")
	return values
}

func main() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on 0.0.0.0:3333")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	bufc := bufio.NewReader(conn)
	defer conn.Close()
	values := getInfo(bufc)
	choice := values[1]
	client := Client {
		conn:  conn,
		name:  values[0],
		chatWith:  "",
	}

	if strings.TrimSpace(client.name) == "" {
		io.WriteString(conn, "Invalid Username\n")
		return
	}

	if choice == "2" {
		value := getInfo(bufc)
		client.chatWith = value[0]
	}

	connections = append(connections, client)

	for {
		msg, err := common.ReadMsg(conn)
		if err != nil {
			if err == io.EOF {
				// Close the connection when you're done with it.
				removeConn(conn)
				conn.Close()
				return
			}
			log.Println(err)
			return
		}

		fmt.Printf("Message Received: %s\n", msg)
		if choice == "2" {
			sendIndividual(client, msg)
		} else if choice == "1" {
			broadcast(conn, msg)
		}
	}
}

func removeConn(conn net.Conn) {
	var i int
	for i = range connections {
		if connections[i].conn == conn {
			break
		}
	}
	connections = append(connections[:i], connections[i+1:]...)
}

func broadcast(conn net.Conn, msg string) {
	for i := range connections {
		if connections[i].conn != conn {
			err := common.WriteMsg(connections[i].conn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func sendIndividual(c Client, msg string) {
        for i := range connections {
                if connections[i].name == c.chatWith {
                        err := common.WriteMsg(connections[i].conn, msg)
                        if err != nil {
                                log.Println(err)
                        }
                }
        }
}
