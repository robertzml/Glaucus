package protocol

import (
	"fmt"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
)


// 报文消息接口
// 所有类型的报文均实现改接口
type Message interface {
	// 报文协议解析
	Parse(payload string) (err error)

	// 打印协议内容
	Print(cell TLV)

	Save()
}


// 处理接收的报文
// topic: 主题
// payload: 接收内容
// qos: QoS
func Receive(topic string, payload []byte, qos byte) {
	cell, msg, err := parseType(string(payload[:]))
	if err != nil {
		fmt.Println("catch error in parseType: ", err.Error())
		return
	}

	err = msg.Parse(cell.Value)
	if err != nil {
		fmt.Println("catch error in parse", err.Error())
		return
	}

	msg.Print(cell)


	msg.Save()
}




