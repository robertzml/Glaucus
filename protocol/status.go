package protocol

import (
	"../equipment"
	"fmt"
	"strconv"
)

// 设备状态报文
type StatusMessage struct {
	SerialNumber    	string
	MainboardNumber 	string
	DeviceType			string
	ControllerType		string
	WaterHeaterStatus	equipment.WaterHeater
}

// 解析协议内容
func (msg *StatusMessage) ParseContent(payload string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch runtime panic: %v\n", err)
		}
	}()

	index := 0
	length := len(payload)

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return
		}

		switch tlv.Tag {
		case 0x127:
			msg.SerialNumber = tlv.Value
		case 0x12b:
			msg.MainboardNumber = tlv.Value
		case 0x125:
			msg.DeviceType = tlv.Value
		case 0x12a:
			msg.ControllerType = tlv.Value
		default:
		}

		if tlv.Tag == 0x128 {
			msg.parseWaterHeater(tlv.Value)
		}

		index += tlv.Length + 8
	}
}

/*
打印协议信息
*/
func (msg* StatusMessage) Print(cell TLV) {
	fmt.Printf("Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

/*
解析热水器状态
 */
func (msg *StatusMessage) parseWaterHeater(payload string) {
	index := 0
	length := len(payload)

	wh := new(equipment.WaterHeater)

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return
		}

		switch tlv.Tag {
		case 0x01:
			wh.Power, _ = strconv.Atoi(tlv.Value)
		case 0x03:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			wh.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			wh.OutFlow = int(v) * 10
		case 0x05:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			wh.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			wh.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			wh.ErrorCode = int(v)
		}

		index += tlv.Length + 8
	}
}