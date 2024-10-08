package models

import (
	"encoding/binary"
	"strings"
)

type Message struct {
  message []byte
  commandCount uint32
}

func NewMessage() *Message {
  return &Message{}
}

func (m *Message) addCommand(command string) {
  cmds := strings.Split(command, " ")
  msg := []byte{command}
  length := uint32(len(message))
  count := m.commandCount + 1
  binary.LittleEndian.PutUint32(m.message[0:4], count)

}
