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

	var whOfflineTopic = fmt.Sprintf("equipment/%d/1/+/offline_info", base.DefaultConfig.MqttChannel)
	if err := m.Subscribe(whOfflineTopic, 0, WaterHeaterOfflineHandler); err != nil {
		fmt.Println(err)
		return
	}

	var whStatusTopic = fmt.Sprintf("equipment/%d/1/+/status_info", base.DefaultConfig.MqttChannel)
	if err := m.Subscribe(whStatusTopic, 0, WaterHeaterStatusHandler); err != nil {
		fmt.Println(err)
		return
	}

	wcStatusTopic := fmt.Sprintf("equipment/%d/2/+/status_info", base.DefaultConfig.MqttChannel)
	if err := m.Subscribe(wcStatusTopic, 0, WaterCleanerStatusHandler); err != nil {
		fmt.Println(err)
	}
}
