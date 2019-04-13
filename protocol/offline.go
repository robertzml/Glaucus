package protocol

import (
	"fmt"
	"strings"
)

// 处理离线报文
// topic: 主题
// payload: 接收内容
// qos: QoS
func Offline(topic string, payload []byte, qos byte) {
	kv := strings.Split(topic, "/")
	if len(kv) != 5 {
		fmt.Println("offline topic is wrong.")
		return
	}

	if string(payload[:]) != "offline" {
		fmt.Println("offline payload is wrong.")
		return
	}

	productType := kv[2]
	serialNumber := kv[3]

	if productType == "1" {
		handleWaterHeaterOffline(serialNumber)
	} else {

	}

	fmt.Printf("equipment %s is offline.\n", serialNumber)
}
