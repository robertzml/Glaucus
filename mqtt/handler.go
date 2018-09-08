/*
 *  消息处理
 */

package mqtt

import (
	"fmt"
	"../protocol"
	paho "github.com/eclipse/paho.mqtt.golang"
)

// 默认订阅消息处理方法
var defaultHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("MSG: %s\n", msg.Payload())

	protocol.Receive(msg.Topic(), msg.Payload(), msg.Qos())
}
