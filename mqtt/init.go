package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/glog"
)

var (
	// 全局MQTT 接收连接
	ReceiveMqtt *MQTT

	// 全局MQTT 发送连接
	SendMqtt *MQTT
)

// 当前包名称
const packageName = "mqtt"

// MQTT 结构体
type MQTT struct {
	ClientId string
	Address  string
	client   paho.Client
}

// paho 日志
type MLogger struct {
	Level int
}

func (m MLogger) Println(v ...interface{}) {
	s := fmt.Sprint(v)
	glog.Write(m.Level, packageName, "paho", s)
}

func (m MLogger) Printf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v)
	glog.Write(m.Level, packageName, "paho", s)
}

// 初始化全局MQTT 变量
func InitMQTT() {
	paho.ERROR = MLogger{1}
	paho.CRITICAL = MLogger{0}
	paho.WARN = MLogger{2}
	// paho.DEBUG = MLogger{4}
}
