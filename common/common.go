package common

import (
	"fmt"
	"net"
	"bufio"
)

func WriteMsg(conn net.Conn, msg string) error {
	_, err := fmt.Fprintf(conn, msg+"\n")
	return err
}

func ReadMsg(conn net.Conn) (string, error) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	return message, err
}
