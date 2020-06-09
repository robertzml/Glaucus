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
			glog.Write(1, packageName, "store", fmt.Sprintf("catch runtime panic in store: %v", r))
		}
	}()

	for {
		pak := <- base.MqttStatusCh
		glog.Write(4, packageName, "store", fmt.Sprintf("TOPIC: %s. MQTT status consumer.", pak.Topic))

		cell, version, seq, msg, err := parseType(pak.ProductType, pak.Payload)
		if err != nil {
			glog.Write(1, packageName, "store", "catch error in parseType: " + err.Error())
			continue
		}

		_ = parseMessage(msg, version, seq, cell)

		glog.Write(4, packageName, "store", "store finish.")
	}
}

// 解析协议
// 根据收到的报文，解析出协议头部，确定协议类型
// cell 报文头
// seq 序列号
// msg  报文内容
func parseType(productType int, message string) (cell tlv.TLV, version float64, seq string, msg Message, err error) {
	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "store", fmt.Sprintf("catch runtime panic in parse type: %v", r))
			err = errors.New("parse type error")
		}
	}()

	// read header
	version, seq, payload, err := tlv.ParseHead(message)
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

// 解析协议内容
func parseMessage(msg Message, version float64, seq string, cell tlv.TLV) (err error) {
	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "store", fmt.Sprintf("catch runtime panic in parse message: %v", r))
			err = errors.New("parse message error")
		}
	}()

	data, err := msg.Parse(cell.Value)
	if err != nil {
		glog.Write(1, packageName, "store", "catch error in parse: " + err.Error())
		return err
	}

	pass := msg.Authorize(seq)
	if !pass {
		glog.Write(2, packageName, "store", "authorize failed.")
		return nil
	}

	err = msg.Handle(data, version, seq)
	if err != nil {
		glog.Write(1, packageName, "store", "catch error in handle: " + err.Error())
		return err
	}

	return nil
}