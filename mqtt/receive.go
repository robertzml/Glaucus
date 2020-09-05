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

	receiveClientId = fmt.Sprintf("influx-receive-test")
	ReceiveMqtt.Connect(receiveClientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttPassword, base.DefaultConfig.MqttServerAddress, receiveOnConnect)
}

// 接收自动订阅
var receiveOnConnect paho.OnConnectHandler = func(client paho.Client) {
	glog.Write(4, packageName, "onConnect", fmt.Sprintf("%s connect to mqtt.", receiveClientId))

	var whStatusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", base.DefaultConfig.MqttChannel)
	if err := ReceiveMqtt.Subscribe(whStatusTopic, 0, WaterHeaterStatusHandler); err != nil {
		glog.Write(1, packageName, "OnConnect", err.Error())
		return
	}
}
