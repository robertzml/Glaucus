package protocol

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
	packageName = "protocol"
)

// 处理接收的报文
// 从 channel 中获取数据
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

		pass, err := msg.Authorize()
		if err != nil {
			glog.Write(1, packageName, "store", "catch error in authorize: " + err.Error())
			return
		}
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
