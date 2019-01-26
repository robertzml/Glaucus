package app

import (
	"../mqtt"
	"fmt"
)

func Run() {
	fmt.Println("app is running")
}

// 启动MQTT 服务订阅
func StartMqtt() {
	var channel = 1
	var clientId = fmt.Sprintf("server-chanel-%d", channel)
	var server = "tcp://192.168.0.120:1883"

	m := new(mqtt.MQTT)

	m.Connect(clientId, server)

	var statusTopic = fmt.Sprintf("equipment/%d/water_heater/+/status_info", channel)
	if err := m.Subscribe(statusTopic, 0, mqtt.StatusHandler); err != nil {
		fmt.Println(err)
		return
	}

	var answerTopic = fmt.Sprintf("equipment/%d/water_heater/+/answer_info", channel)
	if err := m.Subscribe(answerTopic, 2, mqtt.AnswerHandler); err != nil {
		fmt.Println(err)
		return
	}
}