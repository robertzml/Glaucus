package protocol

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/equipment"
	"time"
)

// 净水器设备状态报文
type WCStatusMessage struct {
	SerialNumber    string
	MainboardNumber string
	DeviceType      string
	ControllerType  string
}

// 解析协议内容
func (msg *WCStatusMessage) Parse(payload string) (data interface{}, err error) {
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
func (msg *WCStatusMessage) Print(cell TLV) {
	fmt.Printf("Water Cleaner StatusMessage Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

// 安全检查
// 返回: pass 是否通过
func (msg *WCStatusMessage) Authorize() (pass bool, err error) {
	wcs := new(equipment.WaterCleaner)

	if exists := wcs.LoadStatus(msg.SerialNumber); exists {
		if wcs.MainboardNumber != msg.MainboardNumber {
			return false, errors.New("mainboard Number not equal.")
		}
	} else {
		fmt.Println("authorize: new equipment found.")
		return true, nil
	}

	fmt.Println("authorize: pass.")
	return true, nil
}

// 报文后续处理
func (msg *WCStatusMessage) Handle(data interface{}) (err error) {
	switch data.(type) {
	case TLV:
		tlv := data.(TLV)
		if tlv.Tag == 0x128 {
			// 局部更新
			if err = msg.handleWaterCleanerChange(tlv.Value); err != nil {
				return err
			}
			fmt.Println("water cleaner partial update.")

		} else if tlv.Tag == 0x12e {
			// 整体更新
			if err := msg.handleWaterCleanerTotal(tlv.Value); err != nil {
				return err
			}

			fmt.Println("water cleaner total update.")
		}
	}

	return nil
}

// 整体解析净水器状态
func (msg *WCStatusMessage) handleWaterCleanerTotal(payload string) (err error) {
	wcs := new(equipment.WaterCleaner)

	exists := wcs.LoadStatus(msg.SerialNumber)

	wcs.SerialNumber = msg.SerialNumber
	wcs.MainboardNumber = msg.MainboardNumber
	wcs.Logtime = time.Now().Unix()
	wcs.DeviceType = msg.DeviceType
	wcs.ControllerType = msg.ControllerType

	if !exists || wcs.Online == 0 {
	
	}

	return nil
}

// 处理净水器变化状态，并局部更新
func (msg *WCStatusMessage) handleWaterCleanerChange(payload string) (err error) {
	wcs := new(equipment.WaterCleaner)

	exists := wcs.LoadStatus(msg.SerialNumber)
	if !exists {
		fmt.Println("cannot update partial for new equipment.")
		return nil
	}

	wcs.Logtime = time.Now().Unix()

	return nil
}
