package app

import (
	"../equipment"
	"../mqtt"
	"fmt"
)

func Run() {
	fmt.Println("app is running")

	startMqtt()
}

// 启动MQTT 服务订阅
func startMqtt() {
	var channel = 1
	var clientId = fmt.Sprintf("server-chanel-%d", channel)
	var server = "tcp://192.168.0.120:1883"

	m := new(mqtt.MQTT)

	m.Connect(clientId, server)

	var loginTopic = fmt.Sprintf("equipment/%d/1/+/login_info", channel)
	if err := m.Subscribe(loginTopic, 0, mqtt.LoginHandler); err != nil {
		fmt.Println(err)
		return
	}

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
	var w equipment.WaterHeater

	_, err := w.GetStatus("01100101801100e")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", w)
	}
}

