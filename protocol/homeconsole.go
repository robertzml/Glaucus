package protocol

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
)

// 处理接收的报文
// 从 channel 中获取数据
func Store() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("catch runtime panic in data parse and store: %v\n", r)
		}
	}()

	for {
		pak := <- base.MqttStatusCh
		fmt.Println("store consumer.")

		cell, msg, err := parseType(pak.ProductType, pak.Payload)
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

		fmt.Println("store finish.")
	}
}
