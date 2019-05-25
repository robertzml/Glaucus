package receive

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/tlv"
)

const (
	packageName = "receive"
)

// 处理接收的状态消息报文
// 从 channel 中获取数据，并进行存储
func Store() {
	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "store", fmt.Sprintf("catch runtime panic in data parse and store: %v", r))
		}
	}()

	for {
		pak := <- base.MqttStatusCh
		glog.Write(3, packageName, "store", "mqtt status consumer.")

		cell, msg, err := parseType(pak.ProductType, pak.Payload)
		if err != nil {
			glog.Write(1, packageName, "store", "catch error in parseType: " + err.Error())
			return
		}

		data, err := msg.Parse(cell.Value)
		if err != nil {
			glog.Write(1, packageName, "store", "catch error in parse: " + err.Error())
			return
		}

		pass := msg.Authorize()
		if !pass {
			glog.Write(2, packageName, "store", "authorize failed.")
			return
		}

		err = msg.Handle(data)
		if err != nil {
			glog.Write(1, packageName, "store", "catch error in handle: " + err.Error())
			return
		}

		glog.Write(3, packageName, "store", "store finish.")
	}
}

// 解析协议
// 根据收到的报文，解析出协议头部，确定协议类型
// cell 报文头
// msg  报文内容
func parseType(productType int, message string) (cell tlv.TLV, msg Message, err error) {
	// read header
	_, payload, err := tlv.ParseHead(message)
	if err != nil {
		return
	}

	// parse message cell type
	cell, err = tlv.ParseTLV(payload, 0)
	if err != nil {
		return
	}

	// check message length
	length := len(message) - len(tlv.HomeConsoleVersion) - 8 - 8
	if cell.Length != length {
		err = errors.New("message length is not correct")
		return
	}


	switch cell.Tag {
	case 0x14:
		if productType == 1 {
			msg = new(WHStatusMessage)
		} else if productType == 2 {
			// msg = new(WCStatusMessage)
			msg = nil
			err = errors.New("wrong device type")
		} else {
			msg = nil
			err = errors.New("wrong device type")
		}
	case 0x15:
		if productType == 1 {
			msg = new(WHOfflineMessage)
		}
	default:
		msg = nil
		err = errors.New("TLV not defined")
	}

	return
}
