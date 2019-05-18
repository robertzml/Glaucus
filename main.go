package main

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/mqtt"
	"github.com/robertzml/Glaucus/protocol"
	"github.com/robertzml/Glaucus/redis"
	"github.com/robertzml/Glaucus/rest"
	"time"
)

func main() {
	fmt.Println("app is running")

	defer func() {
		fmt.Println("app is stop.")
	}()

	base.InitConfig()
	base.InitChannel()

	glog.InitGlog()
	go startLog()

	redis.InitPool(base.DefaultConfig.RedisDatabase)

	mqtt.InitMQTT()
	mqtt.InitSend()

	startMqtt()

	go startStore()
	go startRest()
	go startControl()

	for {
		text := fmt.Sprintf("redis active: %d, redis idle: %d. receive mqtt connection: %t, send mqtt connection: %t.",
			redis.RedisPool.ActiveCount(), redis.RedisPool.IdleCount(), mqtt.ReceiveMqtt.IsConnect(), mqtt.SendMqtt.IsConnect())

		glog.Write(3, "main", "state", text)

		time.Sleep(10 * 1e9)
	}
}

// 启动日志服务
func startLog() {
	fmt.Println("start log service.")
	glog.Read()
}

// 启动MQTT 服务订阅
func startMqtt() {
	glog.Write(3, "main", "start", "start mqtt listen.")
	mqtt.StartReceive()
}

// 启用数据存储服务
func startStore() {
	glog.Write(3, "main", "start", "start data store.")
	protocol.Store()
}

// 启动HTTP接收服务
func startRest() {
	glog.Write(3, "main", "start", "start rest server.")
	rest.StartHttpServer()
}

// 启动控制服务
func startControl() {
	glog.Write(3, "main", "start", "start control server.")
	mqtt.StartSend()
}
