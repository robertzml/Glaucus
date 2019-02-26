package protocol

import (
	"errors"
	"fmt"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
)

/*
报文消息接口

所有类型的报文均实现改接口
 */
type Message interface {
	ParseContent(payload string) (err error)
	Print(cell TLV)
	Save()
}

/*
处理接收的报文
 */
func Receive(topic string, payload []byte, qos byte) {
	cell, msg, err := Parse(string(payload[:]))
	if err != nil {
		fmt.Printf("catch error in parse: ", err.Error())
		return
	}

	msg.Print(cell)

	msg.Save()
}


/*
协议解析
根据收到的报文，解析出协议内容
cell 报文头
msg  报文内容
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

	err = msg.ParseContent(payload[8:])

	return
}


