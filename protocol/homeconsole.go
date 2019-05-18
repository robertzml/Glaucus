package protocol

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/mqtt"
	"sync"
	"time"
)

const (
	HomeConsoleVersion = "Homeconsole05.00"
)

// 处理接收的报文
// 从 channel 中获取数据
func Store() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("catch runtime panic in data parse and store: %v\n", r)
		}
	}()

	for {
		pak := <- base.MqttStatusCh
		fmt.Println("store consumer.")

		cell, msg, err := parseType(pak.ProductType, pak.Payload)
		if err != nil {
			fmt.Println("catch error in parseType: ", err.Error())
			return
		}

		data, err := msg.Parse(cell.Value)
		if err != nil {
			fmt.Println("catch error in parse.", err.Error())
			return
		}

		pass, err := msg.Authorize()
		if err != nil {
			fmt.Println("catch error in authorize.", err.Error())
			return
		}
		if !pass {
			fmt.Println("authorize failed.")
			return
		}

		err = msg.Handle(data)
		if err != nil {
			fmt.Println("catch error in handle.", err.Error())
			return
		}

		fmt.Println("store finish.")
	}
}

// 校时任务
func LoopTiming() {
	var wg sync.WaitGroup
	wg.Add(1)
	ticker1 := time.NewTicker(5 * time.Second)

	go func(t *time.Ticker) {
		defer wg.Done()
		for {
			<-t.C
			timing := new(TimingMessage)
			msg := timing.splice()

			var timingTopic = fmt.Sprintf("server/%d/1/+/control_info", base.DefaultConfig.MqttChannel)
			mqtt.SendMqtt.Publish(timingTopic, 1, msg)

			fmt.Println("send timing", msg)
		}
	}(ticker1)

	wg.Wait()
}
