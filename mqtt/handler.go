/*
 *  消息处理
 */

package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
)

// 默认订阅消息处理方法
func defaultHandler(client paho.Client, msg paho.Message) {
	glog.Write(3, packageName, "defaultHandler", fmt.Sprintln("TOPIC: %s, Id: %d, QoS: %d\tMSG: %s",
		msg.Topic(), msg.MessageID(), msg.Qos(), msg.Payload()))
}

// 热水器状态消息订阅处理方法
var WaterHeaterStatusHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	glog.Write(3, packageName, "whstatus",
		fmt.Sprintf("TOPIC: %s, Id: %d, QoS: %d\tMSG: %s", msg.Topic(), msg.MessageID(), msg.Qos(), msg.Payload()))

	pak := new(base.ReceivePacket)
	pak.ProductType = 1
	pak.Topic = msg.Topic()
	pak.Payload = string(msg.Payload()[:])

	base.MqttStatusCh <- pak
	glog.Write(3, packageName, "whstatus", "mqtt status producer.")
}

// 净水器状态消息订阅处理方法
var WaterCleanerStatusHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("Status TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("Status MSG: %s\n", msg.Payload())

	// protocol.Receive(2, msg.Topic(), msg.Payload(), msg.Qos())
}