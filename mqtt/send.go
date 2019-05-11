package mqtt

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
)

// 启动MQTT发送服务
// 通过全局 MqttControlCh 获取发送请求
func StartSend() {
	m := new(MQTT)

	clientId := fmt.Sprintf("send-channel-%d", base.DefaultConfig.MqttChannel)
	m.Connect(clientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttServerAddress)

	defer func() {
		m.Disconnect()
	}()

	for {
		input := <-base.MqttControlCh
		fmt.Println("control consumer.")

		var controlTopic = fmt.Sprintf("server/%d/1/%s/control_info", base.DefaultConfig.MqttChannel, input.SerialNumber)
		m.Publish(controlTopic, 2, input.Payload)

		fmt.Printf("PUBLISH Topic:%s, Payload: %s\n", controlTopic, input.Payload)
	}
}
