package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"go.uber.org/zap"
)

type Message []byte

var (
	suggar *zap.SugaredLogger
)

const (
	BufferSize = 1024
	Port       = 6379
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error initializing logger")
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		panic("Error resolving TCP address")
	}

	suggar = logger.Sugar()

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic("Error connecting")
	}
	defer conn.Close()

	command := "SET test 1"
	message := buildMessage(command)
  sendMessage(conn, message)

	command = "GET test 1"
	message = buildMessage(command)
  sendMessage(conn, message)

}

func buildMessage(command string) Message {
	cmds := strings.Split(command, " ")
	counter := 0
  message := make([]byte, 4)
	binary.LittleEndian.PutUint32(message[0:4], uint32(counter))
	for _, cmd := range cmds {
    suggar.Infof("Appending %s\n", cmd)
		msg := []byte(cmd)

    message = appendUint32(message, uint32(len(msg)))

    message = append(message, msg...)

		counter++
	}
	binary.LittleEndian.PutUint32(message[0:4], uint32(counter))
	return message
}

func sendMessage(conn *net.TCPConn, message Message) {
	fmt.Println(message)
	size, err := conn.Write(message)
	if err != nil {
		panic(err)
	}

  buff := make([]byte, 1024)
  n, err := conn.Read(buff)
  if err != nil {
    suggar.Errorf("Error reading message %v", err)
  }
	suggar.Infof("%d written\n", size)
  suggar.Infof("Reply: %s\n", buff[:n])
}

func appendUint32(message Message, length uint32) []byte {
  var buf [4]byte
  binary.LittleEndian.PutUint32(buf[:], length)
  return append(message, buf[:]...)
}
