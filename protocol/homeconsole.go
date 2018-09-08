package protocol

import "errors"

const (
	HomeConsoleVersion = "Homeconsole02.01"
)

type Message interface {
	ParseContent(payload string)
}


/*
协议解析
根据收到的报文，解析出协议内容
 */
func Parse(message string) (err error) {

	// read header
	_, payload, err := parseHead(message)
	if err != nil {
		return
	}

	// parse message cell type
	cell, err := parseCell(payload)
	if err != nil {
		return
	}

	var msg Message
	switch cell.Tag {
	case 0x14:
		m := new(StatusMessage)
		msg = Message(m)
	default:
		err = errors.New("TLV not defined")
	}

	msg.ParseContent(payload[8:])

	return
}


