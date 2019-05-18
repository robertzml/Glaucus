package mqtt

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	paho "github.com/eclipse/paho.mqtt.golang"
)

var SendMqtt *MQTT

// 初始化全局MQTT连接
func InitSend() {
	SendMqtt = new(MQTT)

	clientId := fmt.Sprintf("send-channel-%d", base.DefaultConfig.MqttChannel)
	SendMqtt.Connect(clientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttServerAddress, sendOnConnect)
}

// 启动MQTT发送服务
// 通过全局 MqttControlCh 获取发送请求
func StartSend() {

	defer func() {
		SendMqtt.Disconnect()
		fmt.Println("Send mqtt function is close.")
	}()

	for {
		input := <-base.MqttControlCh
		fmt.Println("control consumer.")

		var controlTopic = fmt.Sprintf("server/%d/1/%s/control_info", base.DefaultConfig.MqttChannel, input.SerialNumber)
		SendMqtt.Publish(controlTopic, 2, input.Payload)

		fmt.Printf("PUBLISH Topic:%s, Payload: %s\n", controlTopic, input.Payload)
	}
}

var sendOnConnect paho.OnConnectHandler = func(client paho.Client) {
	fmt.Println("connect to mqtt send.")
}