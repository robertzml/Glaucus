package protocol

import (
	"errors"
	"fmt"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
)

type Message interface {
	ParseContent(payload string)
}

/*
处理接收的报文
 */
func Receive(topic string, payload []byte, qos byte) {
	cell, _, err := Parse(string(payload[:]))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Tag is: %d.\n", cell.Tag)
}


/*
协议解析
根据收到的报文，解析出协议内容
 */
func Parse(message string) (cell TLV, msg Message, err error) {

	// read header
	_, payload, err := parseHead(message)
	if err != nil {
		return
	}

	// parse message cell type
	cell, err = parseCell(payload)
	if err != nil {
		return
	}

	switch cell.Tag {
	case 0x03:
		msg = new(LoginMessage)
	case 0x14:
		msg = new(StatusMessage)
	default:
		msg = nil
		err = errors.New("TLV not defined")
	}

	msg.ParseContent(payload[8:])

	return
}


