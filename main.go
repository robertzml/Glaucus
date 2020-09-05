package main

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/influx"
	"github.com/robertzml/Glaucus/mqtt"
	"github.com/robertzml/Glaucus/receive"
	"time"
)

func main() {
	fmt.Println("app is running")

	defer func() {
		fmt.Println("app is stop.")
	}()

	base.LoadConfig(1)
	base.InitChannel()

	glog.InitGlog()
	go startLog()

	influx.InitFlux()
	go startInflux()

	mqtt.InitMQTT()

	// 启动 MQTT订阅，数据处理服务，设备控制服务
	startMqtt()
	go startStore()

	for {
		text := fmt.Sprintf("time is %v.", time.Now())
		glog.Write(4, "main", "state", text)

		time.Sleep(300 * 1e9)
	}
}

// 启动日志服务
func startLog() {
	fmt.Println("start log service.")
	glog.Read()
}

// 启动influx 服务
func startInflux() {
	fmt.Println("start influx service")
	influx.Read()
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

