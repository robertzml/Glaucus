/*
 *  消息处理
 */

package mqtt

import (
	"../protocol"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
)

// 默认订阅消息处理方法
func defaultHandler(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var protocolHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("MSG: %s\n", msg.Payload())

	protocol.Receive(msg.Topic(), msg.Payload(), msg.Qos())
}

// 登录消息订阅处理方法
var LoginHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("Login TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("Login MSG: %s\n", msg.Payload())

	protocol.Receive(msg.Topic(), msg.Payload(), msg.Qos())
}


// 状态消息订阅处理方法
var StatusHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("Status TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("Status MSG: %s\n", msg.Payload())

	protocol.Receive(msg.Topic(), msg.Payload(), msg.Qos())
}

// 响应消息订阅处理方法
var AnswerHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("Answer TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("Answer MSG: %s\n", msg.Payload())
}