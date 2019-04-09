package protocol

import (
	"fmt"
	"github.com/robertzml/Glaucus/equipment"
	"strings"
	"time"
)

// 处理离线报文
// topic: 主题
// payload: 接收内容
// qos: QoS
func Offline(topic string, payload []byte, qos byte) {
	kv := strings.Split(topic, "/")
	if len(kv) != 5 {
		fmt.Println("offline topic is wrong.")
		return
	}

	if string(payload[:]) != "offline" {
		fmt.Println("offline payload is wrong.")
		return
	}

	serialNumber := kv[3]
	whs := new(equipment.WaterHeater)

	if exists := whs.LoadStatus(serialNumber); !exists {
		fmt.Println("don't find equipment.")
		return
	}

	// 更新离线状态和时间
	whs.Online = 0
	whs.LineTime = time.Now().Unix()

	whs.SaveStatus()

	// 关键数据
	whKey := new(equipment.WaterHeaterKey)
	whKey.SerialNumber = whs.SerialNumber
	whKey.MainboardNumber = whs.MainboardNumber
	whKey.Logtime = whs.Logtime
	whKey.Activate = whs.Activate
	whKey.ActivationTime = whs.ActivationTime
	whKey.Lock = whs.Lock
	whKey.DeadlineTime = whs.DeadlineTime
	whKey.Online = 0
	whKey.LineTime = whs.LineTime

	whs.PushKey(whKey)

	fmt.Printf("equipment %s is offline.", serialNumber)
}
