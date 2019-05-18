package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
)

// 全局MQTT 接收连接
var ReceiveMqtt *MQTT

// 全局MQTT 发送连接
var SendMqtt *MQTT


type MQTT struct {
	ClientId string
	Address  string
	client   paho.Client
}

type MLogger struct {

}

func (MLogger) Println(v ...interface{}) {
	fmt.Println(v)
}

func (MLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v)
}

// 初始化全局MQTT 变量
func InitMQTT() {
	paho.ERROR = MLogger{}
	paho.CRITICAL = MLogger{}
	paho.WARN = MLogger{}
	// paho.DEBUG = MLogger{}
}
