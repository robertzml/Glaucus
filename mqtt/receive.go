package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
)

// 启动MQTT接收服务
func StartReceive() {
	ReceiveMqtt = new(MQTT)

	clientId := fmt.Sprintf("receive-channel-%d", base.DefaultConfig.MqttChannel)
	ReceiveMqtt.Connect(clientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttServerAddress, receiveOnConnect)
}

// 接收自动订阅
var receiveOnConnect paho.OnConnectHandler = func(client paho.Client) {
	glog.Write(3, packageName, "onConnect", "receive connect to mqtt.")

	var whStatusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", base.DefaultConfig.MqttChannel)
	if err := ReceiveMqtt.Subscribe(whStatusTopic, 0, WaterHeaterStatusHandler); err != nil {
		glog.Write(1, packageName, "OnConnect", err.Error())
		return
	}

	/*
		wcStatusTopic := fmt.Sprintf("equipment/%d/2/+/status_info", base.DefaultConfig.MqttChannel)
		if err := m.Subscribe(wcStatusTopic, 0, WaterCleanerStatusHandler); err != nil {
			fmt.Println(err)
		}
	*/
}
