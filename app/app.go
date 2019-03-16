package app

import (
	"../mqtt"
	"../redis"
	"../protocol"
	"fmt"
)

func Run() {
	fmt.Println("app is running")

	//startRedis()
	//startMqtt()
	startTest()
}

// 启动redis 线程池
func startRedis() {
	redis.InitPool()
}

// 启动MQTT 服务订阅
func startMqtt() {
	var channel = 1
	var clientId = fmt.Sprintf("server-chanel-%d", channel)
	var server = "tcp://192.168.0.120:1883"

	m := new(mqtt.MQTT)

	m.Connect(clientId, server)

	var statusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", channel)
	if err := m.Subscribe(statusTopic, 0, mqtt.StatusHandler); err != nil {
		fmt.Println(err)
		return
	}

	var answerTopic = fmt.Sprintf("equipment/%d/1/+/answer_info", channel)
	if err := m.Subscribe(answerTopic, 2, mqtt.AnswerHandler); err != nil {
		fmt.Println(err)
		return
	}
}

func startTest() {

	var channel = 1
	var clientId = fmt.Sprintf("server-control-%d", channel)
	var server = "tcp://192.168.0.120:1883"

	m := new(mqtt.MQTT)

	m.Connect(clientId, server)

	serialNumber := "01100101801100e3"

	cm := new(protocol.ControlMessage)
	cm.SerialNumber = serialNumber
	cm.MainboardNumber = "10000000000063"

	payload := cm.Power(1)

	var controlTopic = fmt.Sprintf("server/%d/1/%s/control_info", channel, serialNumber)
	m.Publish(controlTopic, 2, payload)

	fmt.Printf("publish: %s", payload)
}

