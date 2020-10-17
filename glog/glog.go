package glog

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/streadway/amqp"
	"time"
)

const (
	// 日志显示系统名称
	systemName = "Glaucus"
)

// 日志 channel
var logChan  chan *packet

// 日志数据包
type packet struct {
	// 日志级别 0-5
	Level  		int

	// 系统名称
	System 		string

	// 模块名称
	Module		string

	// 操作名称
	Action		string

	// 日志内容
	Message		string
}

// 初始化日志
func InitGlog() {
	logChan = make(chan *packet, 10)
}

// 写日志到channel 中
// {"exception", "error", "waring", "info", "debug", "verbose"}
func Write(level int, module string, action string, message string) {
	pak := packet{Level: level, System: systemName, Module: module, Action: action, Message: message}
	logChan <- &pak
}

// 从channel 中获取日志并写入到队列
func Read() {
	rmConnection, err := amqp.Dial(base.DefaultConfig.RabbitMQAddress)
	if err != nil {
		panic(err)
	}

	rbChannel, err := rmConnection.Channel()
	if err != nil {
		panic(err)
	}

	defer func() {
		rmConnection.Close()
		rbChannel.Close()
		fmt.Println("log service is close.")
	}()

	queue, err := rbChannel.QueueDeclare("LogQueue", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	levels := [...]string{"exception", "error", "waring", "info", "debug", "verbose"}

	for {
		pak := <- logChan

		if pak.Level > base.DefaultConfig.LogLevel {
			continue
		}

		// 获取日志消息内容
		jsonData, _ := json.Marshal(pak)
		// fmt.Println(string(jsonData))

		// 推送到 rabbitmq
		err = rbChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: jsonData,
		})
		if err != nil {
			print(err)
		}

		// 输出到控制台
		if base.DefaultConfig.LogToConsole {
			now := time.Now()
			text := fmt.Sprintf("[%s][%s]-[%s]:[%s]\t%s\n",
				levels[pak.Level], now.Format("2006-01-02 15:04:05.000"), pak.Module, pak.Action, pak.Message)

			fmt.Print(text)
		}
	}
}

