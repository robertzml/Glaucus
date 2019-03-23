package mqtt

import (
	"../base"
	"fmt"
)

// 启动MQTT发送服务
// 通过ch 获取发送请求
func StartSend(ch chan *base.SendPacket) {
	var channel = 1
	var clientId = fmt.Sprintf("send-channel-%d", channel)
	var server = "tcp://192.168.0.120:1883"

	m := new(MQTT)
	m.Connect(clientId, server)

	defer func() {
		m.Disconnect()
	}()

	for {
		input := <- ch
		fmt.Println("control consumer.")
		var controlTopic = fmt.Sprintf("server/%d/1/%s/control_info", channel, input.SerialNumber)
		m.Publish(controlTopic, 2, input.Payload)
		fmt.Printf("PUBLISH  Topic:%s, Payload: %s\n", controlTopic, input.Payload)
	}
}