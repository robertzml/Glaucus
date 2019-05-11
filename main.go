package main

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
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

	redis.InitPool(base.DefaultConfig.RedisDatabase)
	mqtt.InitReceive()

	go startMqtt()
	go startStore()
	go startRest()
	go startControl()

	// go reportStatus()

	for {
		fmt.Printf("time: %s, redis active: %d, redis idle: %d. mqtt connection: %t.\n",
			time.Now(), redis.RedisPool.ActiveCount(), redis.RedisPool.IdleCount(), mqtt.ReceiveMqtt.IsConnect())

		time.Sleep(10 * 1e9)
	}
}

// 启动MQTT 服务订阅
func startMqtt() {
	fmt.Println("start mqtt listen.")
	mqtt.StartReceive()
}

// 启用数据存储服务
func startStore() {
	fmt.Println("start data store.")
	protocol.Store()
}

// 启动HTTP接收服务
func startRest() {
	fmt.Println("start rest server.")
	rest.StartHttpServer()
}

// 启动控制服务
func startControl() {
	fmt.Println("start control server.")
	mqtt.StartSend()
}

func reportStatus() {
	defer func() {
		fmt.Println("report status closed.")
	}()

	for {
		fmt.Printf("time: %s, redis active: %d, redis idle: %d. mqtt connection: %t.\n",
			time.Now(), redis.RedisPool.ActiveCount(), redis.RedisPool.IdleCount(), mqtt.ReceiveMqtt.IsConnect())

		time.Sleep(10 * 1e9)
	}
}
