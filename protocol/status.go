package protocol

import (
	"../equipment"
	"errors"
	"fmt"
	"strconv"
)

// 设备状态报文
type StatusMessage struct {
	SerialNumber    	string
	MainboardNumber 	string
	DeviceType			string
	ControllerType		string
}

// 解析协议内容
func (msg *StatusMessage) Parse(payload string) (data interface{}, err error) {
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
			return nil, err
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
			return tlv, nil
		} else if tlv.Tag == 0x12e {
			return tlv, nil
		}

		index += tlv.Length + 8
	}

	return
}

// 打印协议信息
func (msg *StatusMessage) Print(cell TLV) {
	fmt.Printf("StatusMessage Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

// 安全检查
func (msg *StatusMessage) Authorize() (pass bool, err error) {
	equip := new(equipment.WaterHeater)

	exists, err := equip.GetStatus(msg.SerialNumber)
	if err != nil {
		return false, err
	}

	if !exists {
		fmt.Println("new equipment found.")
		return true, nil
	} else {
		if equip.MainboardNumber != msg.MainboardNumber {
			return false, errors.New("Mainboard Number not equal.")
		}
	}

	fmt.Println("authorize pass.")
	return true, nil
}

// 报文后续处理
func (msg *StatusMessage) Handle(data interface{}) (err error) {
	switch data.(type) {
	case TLV:
		tlv := data.(TLV)
		if tlv.Tag == 0x128 {
			// 局部更新
			err = msg.handleWaterHeaterChange(tlv.Value)
			fmt.Println("partial update.")
			if err != nil {
				return err
			}
		} else if tlv.Tag == 0x12e {
			// 整体更新
			waterHeaterStatus, err := msg.parseWaterHeater(tlv.Value)

			if err != nil {
				return err
			}

			waterHeaterStatus.SaveStatus()
			fmt.Println("total update.")
		}
	}

	return nil
}


/*
解析热水器状态
 */
func (msg *StatusMessage) parseWaterHeater(payload string) (waterHeaterStatus equipment.WaterHeater, err error) {
	index := 0
	length := len(payload)

	waterHeaterStatus.SerialNumber = msg.SerialNumber
	waterHeaterStatus.MainboardNumber = msg.MainboardNumber

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return waterHeaterStatus, err
		}

		switch tlv.Tag {
		case 0x01:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.Power = int8(v)
		case 0x03:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.OutFlow = int(v) * 10
		case 0x05:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.ErrorCode = int(v)
		case 0x08:
			waterHeaterStatus.WifiVersion = tlv.Value
		case 0x09:
			v, _ := ParseTime(tlv.Value)
			waterHeaterStatus.CumulateHeatTime = v
		case 0x0a:
			v, _ := ParseCumulate(tlv.Value, 8)
			waterHeaterStatus.CumulateHotWater = v
		case 0x0b:
			v, _ := ParseTime(tlv.Value)
			waterHeaterStatus.CumulateWorkTime = v
		case 0x0c:
			v, _ := ParseCumulate(tlv.Value, 8)
			waterHeaterStatus.CumulateUsedPower = v
		case 0x0d:
			v, _ := ParseCumulate(tlv.Value, 8)
			waterHeaterStatus.CumualteSavePower = v
		case 0x1a:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.Lock = int8(v)
		case 0x1b:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.Activate = int8(v)
		case 0x1c:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.SetTemp = int(v)
		case 0x1d:
			waterHeaterStatus.SoftwareFunction = tlv.Value
		case 0x1e:
			v, _ := ParseCumulate(tlv.Value, 4)
			waterHeaterStatus.OutputPower = v
		case 0x1f:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.ManualClean = int8(v)
		case 0x20:
			v, _ := ParseDateToTimestamp(tlv.Value)
			waterHeaterStatus.DeadlineTime = v
		case 0x21:
			v, _ := ParseDateToTimestamp(tlv.Value)
			waterHeaterStatus.ActivationTime = v
		case 0x22:
			waterHeaterStatus.SpecialParameter = tlv.Value
		}

		index += tlv.Length + 8
	}

	return
}


// 处理热水器变化状态，并局部更新
func (msg *StatusMessage) handleWaterHeaterChange(payload string) (err error) {
	var whs equipment.WaterHeater

	exists, err := whs.GetStatus(msg.SerialNumber)
	if err != nil {
		return err
	}

	if !exists {
		fmt.Println("cannot update partial for new equipment.")
		return nil
	}

	index := 0
	length := len(payload)

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return err
		}

		switch tlv.Tag {
		case 0x01:
			v, _ := strconv.Atoi(tlv.Value)
			whs.Power = int8(v)
		case 0x03:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.OutFlow = int(v) * 10
		case 0x05:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.ErrorCode = int(v)
		case 0x08:
			whs.WifiVersion = tlv.Value
		case 0x09:
			v, _ := ParseTime(tlv.Value)
			whs.CumulateHeatTime = v
		case 0x0a:
			v, _ := ParseCumulate(tlv.Value, 8)
			whs.CumulateHotWater = v
		case 0x0b:
			v, _ := ParseTime(tlv.Value)
			whs.CumulateWorkTime = v
		case 0x0c:
			v, _ := ParseCumulate(tlv.Value, 8)
			whs.CumulateUsedPower = v
		case 0x0d:
			v, _ := ParseCumulate(tlv.Value, 8)
			whs.CumualteSavePower = v
		case 0x1a:
			v, _ := strconv.Atoi(tlv.Value)
			whs.Lock = int8(v)
		case 0x1b:
			v, _ := strconv.Atoi(tlv.Value)
			whs.Activate = int8(v)
		case 0x1c:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.SetTemp = int(v)
		case 0x1d:
			whs.SoftwareFunction = tlv.Value
		case 0x1e:
			v, _ := ParseCumulate(tlv.Value, 4)
			whs.OutputPower = v
		case 0x1f:
			v, _ := strconv.Atoi(tlv.Value)
			whs.ManualClean = int8(v)
		case 0x20:
			v, _ := ParseDateToTimestamp(tlv.Value)
			whs.DeadlineTime = v
		case 0x21:
			v, _ := ParseDateToTimestamp(tlv.Value)
			whs.ActivationTime = v
		case 0x22:
			whs.SpecialParameter = tlv.Value
		}

		index += tlv.Length + 8
	}

	whs.SaveStatus()

	return nil
}