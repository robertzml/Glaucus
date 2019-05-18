package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/base"
)

// 全局MQTT 接收连接
var ReceiveMqtt *MQTT

// 初始化全局MQTT连接
func InitReceive() {
	paho.ERROR = MLogger{}
	paho.CRITICAL = MLogger{}
	paho.WARN = MLogger{}
	// paho.DEBUG = MLogger{}

	ReceiveMqtt = new(MQTT)
}


// 启动MQTT接收服务
func StartReceive() {

	clientId := fmt.Sprintf("receive-channel-%d", base.DefaultConfig.MqttChannel)
	ReceiveMqtt.Connect(clientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttServerAddress, receiveOnConnect)

}

// 接收自动订阅
var receiveOnConnect paho.OnConnectHandler = func(client paho.Client) {
	var whStatusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", base.DefaultConfig.MqttChannel)
	if err := ReceiveMqtt.Subscribe(whStatusTopic, 0, WaterHeaterStatusHandler); err != nil {
		fmt.Println(err)
		return
	}

	/*
		wcStatusTopic := fmt.Sprintf("equipment/%d/2/+/status_info", base.DefaultConfig.MqttChannel)
		if err := m.Subscribe(wcStatusTopic, 0, WaterCleanerStatusHandler); err != nil {
			fmt.Println(err)
		}
	*/
}
