package glog

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"io"
	"os"
	"path/filepath"
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
	folder, _ := filepath.Abs("./log")
	createDir(folder)

	logChan = make(chan *packet, 10)
}

// 异常日志
func WriteException(module string, action string, message string) {
	Write(0, module, action, message)
}

// 错误日志
func WriteError(module string, action string, message string) {
	Write(1, module, action, message)
}

// 警告日志
func WriteWarning(module string, action string, message string) {
	Write(2, module, action, message)
}

// 信息日志
func WriteInfo(module string, action string, message string) {
	Write(3, module, action, message)
}

// 调试日志
func WriteDebug(module string, action string, message string) {
	Write(4, module, action, message)
}

// 冗余日志
func WriteVerbose(module string, action string, message string) {
	Write(5, module, action, message)
}

// 写日志到channel 中
// {"exception", "error", "waring", "info", "debug", "verbose"}
func Write(level int, module string, action string, message string) {
	pak := packet{Level: level, System: systemName, Module: module, Action: action, Message: message}
	logChan <- &pak
}

// 从channel 中获取日志并写入到队列
func Read() {
	//rmConnection, err := amqp.Dial(base.DefaultConfig.RabbitMQAddress)
	//if err != nil {
	//	panic(err)
	//}
	//
	//rbChannel, err := rmConnection.Channel()
	//if err != nil {
	//	panic(err)
	//}

	defer func() {
		//rmConnection.Close()
		//rbChannel.Close()
		fmt.Println("log service is close.")
	}()

	//queue, err := rbChannel.QueueDeclare("LogQueue", true, false, false, false, nil)
	//if err != nil {
	//	panic(err)
	//}

	levels := [...]string{"exception", "error", "waring", "info", "debug", "verbose"}

	for {
		pak := <- logChan

		if pak.Level > base.DefaultConfig.LogLevel {
			continue
		}

		// 获取日志消息内容
		//jsonData, _ := json.Marshal(pak)
		// fmt.Println(string(jsonData))

		// 推送到 rabbitmq
		//err = rbChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		//	DeliveryMode: amqp.Persistent,
		//	ContentType: "text/plain",
		//	Body: jsonData,
		//})
		//if err != nil {
		//	print(err)
		//}

		now := time.Now()
		filename := fmt.Sprintf("./log/%d%02d%02d.log", now.Year(), now.Month(), now.Day())
		path, _ := filepath.Abs(filename)

		text := fmt.Sprintf("[%s][%s]-[%s]:[%s][%s]\t%s\n",
			levels[pak.Level], now.Format("2006-01-02 15:04:05.000"), pak.System, pak.Module, pak.Action, pak.Message)

		if err := writeFile(path, []byte(text)); err != nil {
			fmt.Println(err)
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

// 创建文件夹
func createDir(path string) {
	_, err := os.Stat(path)
	if err != nil{
		if os.IsNotExist(err){
			_ = os.Mkdir(path, 0744)
		}
	}
}

// 写文件
func writeFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0x644)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return err
}
