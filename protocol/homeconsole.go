package protocol

import (
	"fmt"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
)


// 处理接收的报文
// deviceType: 设备类型
// topic: 主题
// payload: 接收内容
// qos: QoS
func Receive(deviceType int, topic string, payload []byte, qos byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("catch runtime panic in mqtt receive: %v\n", r)
		}
	}()

	cell, msg, err := parseType(deviceType, string(payload[:]))
	if err != nil {
		fmt.Println("catch error in parseType: ", err.Error())
		return
	}

	data, err := msg.Parse(cell.Value)
	if err != nil {
		fmt.Println("catch error in parse.", err.Error())
		return
	}

	pass, err := msg.Authorize()
	if !pass {
		fmt.Println("authorize failed.")
		return
	}
	if err != nil {
		fmt.Println("catch error in authorize.", err.Error())
		return
	}

	err = msg.Handle(data)
	if err != nil {
		fmt.Println("catch error in handle.", err.Error())
		return
	}
}
