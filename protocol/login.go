package protocol

import (
	"fmt"
	"strconv"
)

/*
设备登录报文
 */
type LoginMessage struct {
	SerialNumber    string
	MainboardNumber string
	DeviceType	int
	ControllerType	string
}

func (msg *LoginMessage) ParseContent(payload string) (err error) {
	var index = 0
	length := len(payload)

	for index <= length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return err
		}

		switch tlv.Tag {
		case 0x0127:
			msg.SerialNumber = tlv.Value
		case 0x12b:
			msg.MainboardNumber = tlv.Value
		case 0x125:
			msg.DeviceType, _ = strconv.Atoi(tlv.Value)
		case 0x12a:
			msg.ControllerType = tlv.Value
		default:
		}

		index += tlv.Length + 8
	}

	return
}

func (msg *LoginMessage) Print(cell TLV) {

}

func (msg *LoginMessage) Save() {

}