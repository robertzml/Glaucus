package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
)

var receiveClientId string

// 启动MQTT接收服务
func StartReceive() {
	ReceiveMqtt = new(MQTT)

	receiveClientId = fmt.Sprintf("glaucs-3-test")
	ReceiveMqtt.Connect(receiveClientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttPassword, base.DefaultConfig.MqttServerAddress, receiveOnConnect)
}

// 接收自动订阅
var receiveOnConnect paho.OnConnectHandler = func(client paho.Client) {
	glog.WriteDebug(packageName, "onConnect", fmt.Sprintf("%s connect to mqtt.", receiveClientId))

	// 订阅热水器状态
	var whStatusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", base.DefaultConfig.MqttChannel)
	if err := ReceiveMqtt.Subscribe(whStatusTopic, 0, WaterHeaterStatusHandler); err != nil {
		glog.WriteError(packageName, "OnConnect", err.Error())
		return
	}
}
