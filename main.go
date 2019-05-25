package main

import (
	"flag"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/mqtt"
	"github.com/robertzml/Glaucus/receive"
	"github.com/robertzml/Glaucus/redis"
	"github.com/robertzml/Glaucus/rest"
	"time"
)

var sendMode = flag.Bool("send", false, "启用发送服务")
var receiveMode = flag.Bool("receive", true, "启用接收服务")
var channelId = flag.Int("c", 1, "接收频道")


func main() {
	flag.Parse()

	fmt.Println("app is running")

	defer func() {
		fmt.Println("app is stop.")
	}()


	if *sendMode {
		fmt.Println("in send mode.")

		base.InitConfig(0)
		base.InitChannel()
		glog.InitGlog()
		go startLog()

		redis.InitPool(base.DefaultConfig.RedisDatabase)
		mqtt.InitMQTT()
		mqtt.InitSend()

		// 启动 接口服务，设备控制服务
		go startRest()
		go startControl()

		for {
			text := fmt.Sprintf("redis active: %d, redis idle: %d. send mqtt connection: %t.",
				redis.RedisPool.ActiveCount(), redis.RedisPool.IdleCount(), mqtt.SendMqtt.IsConnect())

			glog.Write(3, "main", "state", text)

			time.Sleep(10 * 1e9)
		}

	} else if *receiveMode {
		fmt.Printf("in receive mode, channel is %d.\n", *channelId)

		base.InitConfig(*channelId)
		base.InitChannel()
		glog.InitGlog()
		go startLog()

		redis.InitPool(base.DefaultConfig.RedisDatabase)
		mqtt.InitMQTT()

		// 启动 MQTT订阅，数据处理服务，设备控制服务
		startMqtt()
		go startStore()
		go startControl()

		for {
			text := fmt.Sprintf("redis active: %d, redis idle: %d. receive mqtt connection: %t, send mqtt connection: %t.",
				redis.RedisPool.ActiveCount(), redis.RedisPool.IdleCount(), mqtt.ReceiveMqtt.IsConnect(), mqtt.SendMqtt.IsConnect())

			glog.Write(3, "main", "state", text)

			time.Sleep(10 * 1e9)
		}
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
	receive.Store()
}

// 启动设备控制服务
func startControl() {
	glog.Write(3, "main", "start", "start control server.")
	mqtt.StartSend()
}

// 启动HTTP接收服务
func startRest() {
	glog.Write(3, "main", "start", "start rest server.")
	rest.StartHttpServer()
}


