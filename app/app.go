package app

import (
	"fmt"
	"github.com/robertzml/Glaucus/protocol"
	"log"
	"time"

	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/mqtt"
	"github.com/robertzml/Glaucus/redis"
	"github.com/robertzml/Glaucus/rest"
)

func Run() {
	fmt.Println("app is running")

	base.InitConfig()
	base.InitChannel()

	startRedis()

	go startMqtt()
	go startStore()
	go startRest()
	go startControl()

	go reportStatus()
	//startTest()
}

// 启动redis 线程池
func startRedis() {
	fmt.Println("start redis pool.")
	redis.InitPool(base.DefaultConfig.RedisDatabase)
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
	for {
		fmt.Printf("time: %s, redis active: %d, redis idle: %d. mqtt connection: %t.\n",
			time.Now(), redis.RedisPool.ActiveCount(), redis.RedisPool.IdleCount(), mqtt.ReceiveMqtt.IsConnect())

		time.Sleep(10 * 1e9)
	}
}

func startTest() {
	mqtt.GetConnections()

	log.Println("abc")
}
