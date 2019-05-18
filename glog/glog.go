package glog

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"io"
	"os"
	"path/filepath"
	"time"
)

// 日志 channel
var GlogCh  chan *GlogPacket


// 日志数据包
type GlogPacket struct {
	// 日志级别
	Level  		int

	// 模块名称
	PackageName	string

	// 标题
	Title		string

	// 日志内容
	Message		string
}

// 初始化日志
func InitGlog() {
	folder, _ := filepath.Abs("./log")
	createDir(folder)

	GlogCh = make(chan *GlogPacket, 10)
}

// 写日志到channel 中
func Write(level int, packageName string, title string, message string) {
	packet := GlogPacket{Level: level, PackageName: packageName, Title: title, Message: message}
	GlogCh <- &packet
}

// 从channel 中获取日志并写入到文件
func Read() {
	defer func() {
		fmt.Println("log service is close.")
	}()

	levels := [...]string{"exception", "error", "waring", "info", "debug"}

	for {
		packet := <- GlogCh

		if packet.Level > base.DefaultConfig.LogLevel {
			continue
		}

		now := time.Now()
		filename := fmt.Sprintf("./log/%d%02d%02d.log", now.Year(), now.Month(), now.Day())
		path, _ := filepath.Abs(filename)

		text := fmt.Sprintf("[%s][%s]-[%s]:[%s]\t%s\n",
			levels[packet.Level], now.Format("2006-01-02 15:04:05"), packet.PackageName, packet.Title, packet.Message)

		if err := writeFile(path, []byte(text)); err != nil {
			fmt.Println(err)
		}
		if base.DefaultConfig.LogToConsole {
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