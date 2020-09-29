package main

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/influx"
	"github.com/robertzml/Glaucus/mqtt"
	"github.com/robertzml/Glaucus/receive"
)

func main() {
	fmt.Println("app is running")

	defer func() {
		fmt.Println("app is stop.")
	}()

	// 载入配置文件
	base.LoadConfig()

	// 初始化全局channel
	base.InitChannel()

	// 启动日志服务
	glog.InitGlog()
	go startLog()

	//influx.InitFlux()
	//go startInflux()

	// 启动 MQTT订阅服务
	mqtt.InitMQTT()
	startMqtt()

	// 启动数据处理服务
	go startProcess()

	// 阻塞
	select{}
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

// 启用数据处理服务
func startProcess() {
	glog.Write(3, "main", "start", "start data process.")
	receive.Process()
}

