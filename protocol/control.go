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

func (msg *ControlMessage) Splice() string {
	head := spliceHead()

	sn := TLV{ Tag: 0x127, Length: len(msg.SerialNumber), Value:msg.SerialNumber }
	mn := TLV{ Tag: 0x12b, Length: len(msg.MainboardNumber), Value:msg.MainboardNumber }
	ca := TLV{ Tag: 0x12, Length: len(msg.ControlAction), Value:msg.ControlAction }

	v := spliceTLV(0x0010, sn.String() + mn.String() + ca.String())

	return head + v
}

// 从缓存中读取设备主板序列号
func (msg *ControlMessage) loadEquipment() bool {
	rc := new(redis.RedisClient)
	rc.Get()
	defer rc.Close()

	mn := rc.Hget("wh_" + msg.SerialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return false
	} else {
		msg.MainboardNumber = mn
		return true
	}
}

func (msg *ControlMessage) Power(power int) string {
	msg.ControlAction = spliceTLV(0x01, strconv.Itoa(power))

	return msg.Splice()
}

func (msg *ControlMessage) Lock(isLock int) string {
	msg.ControlAction = spliceTLV(0x1a, strconv.Itoa(isLock))

	return msg.Splice()
}