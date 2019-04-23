package app

import (
	"fmt"
	"log"

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
	startMqtt()

	go startRest()
	go startControl()

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

func startTest() {
	mqtt.GetConnections()

	log.Println("abc")
}
