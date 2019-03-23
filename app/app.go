package app

import (
	"../base"
	"../mqtt"
	"../redis"
	"../rest"
	"fmt"
)

func Run() {
	fmt.Println("app is running")

	startRedis()
	//startMqtt()
	//startTest()

	ch := make(chan *base.SendPacket)
	go startRest(ch)
	go startControl(ch)
}

// 启动redis 线程池
func startRedis() {
	fmt.Println("start redis pool.")
	redis.InitPool()
}

// 启动MQTT 服务订阅
func startMqtt() {
	fmt.Println("start mqtt listen.")
	var channel = 1
	var clientId = fmt.Sprintf("server-channel-%d", channel)
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

func startRest(ch chan *base.SendPacket) {
	fmt.Println("start rest server.")
	rest.StartHttpServer(ch)
}

// 启动控制服务
func startControl(ch chan *base.SendPacket) {
	mqtt.StartSend(ch)
}

func startTest() {

}

