package main

import (
	"bufio"
	"fmt"
	"chat-server-golang/common"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Connect to server through tcp.
	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go printOutput(conn)
	writeInput(conn)
}

func writeInput(conn net.Conn) {
	fmt.Print("Enter username: ")
	// Read from stdin.
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	username = username[:len(username)-1]
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\nSelect one of the following:\n1. Join Chat Room\t\t2. Individual Chat\n")
	fmt.Print("Press 1 or 2 :: ")
	choice, err := reader.ReadString('\n')
        choice = choice[:len(choice)-1]
        if err != nil {
                log.Fatal(err)
        }

	if choice != "1" && choice != "2" {
		fmt.Println("\nInvalid Choice, Now exiting !!!")
		os.Exit(1)
	}

	err = common.WriteMsg(conn, username+":"+choice)
	if err != nil {
		log.Println(err)
        }

	if choice == "2" {
		fmt.Print("\nWho do you want to chat with: ")
		name, err := reader.ReadString('\n')
	        name = name[:len(name)-1]
		if err != nil {
			log.Fatal(err)
		}
		err = common.WriteMsg(conn, name)
        	if err != nil {
                	log.Println(err)
        	}

	}

	fmt.Println("Enter text: ")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		err = common.WriteMsg(conn, username+": "+text)
		if err != nil {
			log.Println(err)
		}
	}
}

func printOutput(conn net.Conn) {
	for {

		msg, err := common.ReadMsg(conn)
		// Receiving EOF means that the connection has been closed
		if err == io.EOF {
			// Close conn and exit
			conn.Close()
			fmt.Println("Connection Closed. Bye bye.")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(msg)
	}
}
