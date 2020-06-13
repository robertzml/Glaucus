package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
)

var sendClientId string

// 初始化发送
func InitSend() {
	SendMqtt = new(MQTT)
}

// 启动MQTT发送服务
// 通过全局 MqttControlCh 获取发送请求
func StartSend() {
	defer func() {
		SendMqtt.Disconnect()
		fmt.Println("Send mqtt function is close.")
	}()

	sendClientId = fmt.Sprintf("send-channel-%d", base.DefaultConfig.MqttChannel)
	SendMqtt.Connect(sendClientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttPassword, base.DefaultConfig.MqttServerAddress, sendOnConnect)

	for {
		input := <-base.MqttControlCh
		glog.Write(4, packageName, "send", fmt.Sprintf("sn: %s. MQTT control consumer.", input.SerialNumber))

		var controlTopic = fmt.Sprintf("server/1/%s/control_info", input.SerialNumber)
		SendMqtt.Publish(controlTopic, 2, input.Payload)

		glog.Write(4, packageName, "send", fmt.Sprintf("PUBLISH Topic:%s, Payload: %s", controlTopic, input.Payload))
	}
}

// 发送连接回调
var sendOnConnect paho.OnConnectHandler = func(client paho.Client) {
	glog.Write(3, packageName, "onConnect", fmt.Sprintf("%s connect to mqtt.", sendClientId))
}
