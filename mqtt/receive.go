package mqtt

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
)

// 启动MQTT接收服务
func StartReceive() {
	m := new(MQTT)

	clientId := fmt.Sprintf("server-channel-%d", base.DefaultConfig.MqttChannel)
	m.Connect(clientId, base.DefaultConfig.MqttServerAddress)

	var offlineTopic = fmt.Sprintf("equipment/%d/1/+/offline_info", base.DefaultConfig.MqttChannel)
	if err := m.Subscribe(offlineTopic, 0, OfflineHandler); err != nil {
		fmt.Println(err)
		return
	}

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
