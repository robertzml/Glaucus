package send

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/streadway/amqp"
)

/*
下发设备控制指令模块
统一接收控制指令，并传送到rabbitmq
 */

const (
	packageName = "send"
)

// 下发控制 channel
var sendChan  chan *packet

// 下发通信包
type packet struct {
	SerialNumber    string
	MainboardNumber string
	DeviceType		int		// 1:热水器
	ControlType		int		// 控制类型
	Parameter1		string	// 控制参数1
	Parameter2		string 	// 控制参数2
}

/*
初始化下发控制channel
 */
func InitSend() {
	sendChan = make(chan *packet, 10)
}

/*
下发控制指令到channel
 */
func Write(serialNumber string, mainboardNumber string, controlType int, parameter1 string) {
	pak := packet{SerialNumber: serialNumber, MainboardNumber: mainboardNumber, DeviceType: 1,
		ControlType: controlType, Parameter1: parameter1, Parameter2: ""}
	sendChan <- &pak
}

/*
下发控制指令到channel
2个参数
*/
func WriteExtend(serialNumber string, mainboardNumber string, controlType int, parameter1 string, parameter2 string) {
	pak := packet{SerialNumber: serialNumber, MainboardNumber: mainboardNumber, DeviceType: 1,
		ControlType: controlType, Parameter1: parameter1, Parameter2: parameter2}
	sendChan <- &pak
}

/*
从channel 中获取控制指令并写入到队列
 */
func Read(){
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
		fmt.Println("send service is close.")
	}()

	queue, err := rbChannel.QueueDeclare("ControlQueue", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for {
		pak := <- sendChan

		// 获取下发控制指令内容
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
	}
}