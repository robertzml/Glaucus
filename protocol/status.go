package protocol

import (
	"../equipment"
	"../redis"
	"fmt"
	"strconv"
)

// 设备状态报文
type StatusMessage struct {
	SerialNumber    	string
	MainboardNumber 	string
	DeviceType			string
	ControllerType		string
	FullStatus			bool
	WaterHeaterStatus	equipment.WaterHeater
}

// 解析协议内容
func (msg *StatusMessage) ParseContent(payload string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("catch runtime panic: %v\n", r)
			err = fmt.Errorf("%v", r)
		}
	}()

	index := 0
	length := len(payload)

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return err
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
			msg.FullStatus = false
			// msg.parseWaterHeater(tlv.Value)
		} else if tlv.Tag == 0x12e {
			msg.FullStatus = true
			msg.parseWaterHeater(tlv.Value)
		}

		index += tlv.Length + 8
	}

	return
}

/*
打印协议信息
*/
func (msg *StatusMessage) Print(cell TLV) {
	fmt.Printf("StatusMessage Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

/*
保存设备状态信息
 */
func (msg *StatusMessage) Save() {
	r := new(redis.Redis)
	defer r.Close()

	r.Connect()

	if msg.FullStatus {
		r.Hmset(msg.SerialNumber, msg.WaterHeaterStatus)
	}
}


/*
解析热水器状态
 */
func (msg *StatusMessage) parseWaterHeater(payload string) {
	index := 0
	length := len(payload)

	msg.WaterHeaterStatus.MainboardNumber = msg.MainboardNumber

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return
		}

		switch tlv.Tag {
		case 0x01:
			msg.WaterHeaterStatus.Power, _ = strconv.Atoi(tlv.Value)
		case 0x03:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			msg.WaterHeaterStatus.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			msg.WaterHeaterStatus.OutFlow = int(v) * 10
		case 0x05:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			msg.WaterHeaterStatus.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			msg.WaterHeaterStatus.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			msg.WaterHeaterStatus.ErrorCode = int(v)
		}

		index += tlv.Length + 8
	}
}