package protocol

import (
	"fmt"
	"strconv"
)

// 设备状态报文
type StatusMessage struct {
	SerialNumber    string
	MainboardNumber string
	Power byte
	OutputWaterTemp byte
}

// 解析协议内容
func (msg *StatusMessage) ParseContent(payload string) {
	index := 0
	length := len(payload)

	for index <= length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return
		}

		switch tlv.Tag {
		case 0x01:
			msg.SerialNumber = tlv.Value
		case 0x12b:
			msg.MainboardNumber = tlv.Value
		case 0x128:

		}

		index += tlv.Length
	}
}

func (msg *StatusMessage) parseHotHeater(payload string) {
	index := 0
	length := len(payload)

	for index <= length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return
		}

		switch tlv.Tag {
		case 0x01:
			v, _ := strconv.ParseUint(tlv.Value, 16, 0)
			msg.Power = byte(v)
		case 0x03:
			v, _ := strconv.ParseUint(tlv.Value, 16, 0)
			msg.OutputWaterTemp = byte(v)
		}

		index += tlv.Length
	}
}