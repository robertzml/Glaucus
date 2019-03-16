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

	// msg.Print(cell)

	pass, err := msg.Authorize()
	if !pass || err != nil {
		fmt.Println("catch error in authorize.", err.Error())
		return
	}

	err = msg.Handle(data)
	if err != nil {
		fmt.Println("catch error in handle.", err.Error())
		return
	}
}




