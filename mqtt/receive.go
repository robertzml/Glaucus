package mqtt

import (
	"../base"
	"fmt"
)

// 启动MQTT接收服务
func StartReceive() {
	var clientId = fmt.Sprintf("server-channel-%d", base.DefaultConfig.MqttChannel)

	m := new(MQTT)
	m.Connect(clientId, base.DefaultConfig.MqttServerAddress)

	var statusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", base.DefaultConfig.MqttChannel)
	if err := m.Subscribe(statusTopic, 0, StatusHandler); err != nil {
		fmt.Println(err)
		return
	}

	/*
	var answerTopic = fmt.Sprintf("equipment/%d/1/+/answer_info", base.DefaultConfig.MqttChannel)
	if err := m.Subscribe(answerTopic, 2, AnswerHandler); err != nil {
		fmt.Println(err)
		return
	}
	*/
}