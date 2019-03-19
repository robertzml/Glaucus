package protocol

import (
	"../redis"
	"strconv"
)

// 设备控制报文
type ControlMessage struct {
	SerialNumber    	string
	MainboardNumber 	string
	ControlAction		string
}

// 拼接设备控制报文
func (msg *ControlMessage) splice() string {
	head := spliceHead()

	sn := spliceTLV(0x127, msg.SerialNumber)
	mn := spliceTLV(0x12b, msg.MainboardNumber)
	ca := spliceTLV(0x012, msg.ControlAction)

	body := spliceTLV(0x0010, sn + mn + ca)

	return head + body
}

// 从缓存中读取设备主板序列号
// serialNumber: 设备序列号
func (msg *ControlMessage) LoadEquipment(serialNumber string) bool {
	rc := new(redis.RedisClient)
	rc.Get()
	defer rc.Close()

	msg.SerialNumber = serialNumber
	mn := rc.Hget("wh_" + msg.SerialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return false
	} else {
		msg.MainboardNumber = mn
		return true
	}
}

// 开关机报文
func (msg *ControlMessage) Power(power int) string {
	msg.ControlAction = spliceTLV(0x01, strconv.Itoa(power))

	return msg.splice()
}

func (msg *ControlMessage) Lock(isLock int) string {
	msg.ControlAction = spliceTLV(0x1a, strconv.Itoa(isLock))

	return msg.splice()
}