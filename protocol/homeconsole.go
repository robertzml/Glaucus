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
	// data: 返回的数据
	Parse(payload string) (data interface{}, err error)

	// 打印协议内容
	Print(cell TLV)

	// 安全检查
	Authorize() (pass bool, err error)

	// 报文后续处理
	Handle(data interface{}) (err error)
}

// 处理接收的报文
// topic: 主题
// payload: 接收内容
// qos: QoS
func Receive(topic string, payload []byte, qos byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("catch runtime panic in mqtt receive: %v\n", r)
		}
	}()

	cell, msg, err := parseType(string(payload[:]))
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
