package protocol

import (
	"github.com/robertzml/Glaucus/redis"
	"strconv"
)

// 设备控制报文
type WHControlMessage struct {
	SerialNumber    string
	MainboardNumber string
	ControlAction   string
}

// 从缓存中读取设备主板序列号
// serialNumber: 设备序列号
func (msg *WHControlMessage) LoadEquipment(serialNumber string) bool {
	rc := new(redis.RedisClient)
	rc.Get()
	defer rc.Close()

	msg.SerialNumber = serialNumber
	mn := rc.Hget("wh_"+msg.SerialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return false
	} else {
		msg.MainboardNumber = mn
		return true
	}
}

// 开关机报文
func (msg *WHControlMessage) Power(power int) string {
	msg.ControlAction = spliceTLV(0x01, strconv.Itoa(power))
	return msg.splice()
}

// 激活非激活报文
func (msg *WHControlMessage) Activate(status int) string {
	msg.ControlAction = spliceTLV(0x1b, strconv.Itoa(status))
	return msg.splice()
}

// 设备加锁报文
func (msg *WHControlMessage) Lock() string {
	msg.ControlAction = spliceTLV(0x1a, strconv.Itoa(0))
	return msg.splice()
}

// 设备解锁报文
func (msg *WHControlMessage) Unlock(deadline int64) string {
	unlock := spliceTLV(0x1a, strconv.Itoa(1))
	dl := parseTimestampToString(deadline)

	msg.ControlAction = unlock + spliceTLV(0x20, dl)
	return msg.splice()
}

// 设定温度报文
func (msg *WHControlMessage) SetTemp(temp int) string {
	msg.ControlAction = spliceTLV(0x1c, strconv.FormatInt(int64(temp), 16))
	return msg.splice()
}

// 设置允许使用时间
func (msg *WHControlMessage) SetDeadline(deadline int64) string {
	dl := parseTimestampToString(deadline)
	msg.ControlAction = spliceTLV(0x20, dl)
	return msg.splice()
}

// 手动清洗开关
func (msg *WHControlMessage) Clean(status int) string {
	msg.ControlAction = spliceTLV(0x1f, strconv.Itoa(status))
	return msg.splice()
}

// 清零报文
func (msg *WHControlMessage) Clear(status int8) string {
	if status&0x01 == 0x01 {
		msg.ControlAction = spliceTLV(0x39, strconv.Itoa(0))
	} else if status&0x02 == 0x02 {
		msg.ControlAction = spliceTLV(0x38, strconv.Itoa(0))
	} else if status&0x04 == 0x04 {
		msg.ControlAction = spliceTLV(0x37, strconv.Itoa(0))
	} else if status&0x08 == 0x08 {
		msg.ControlAction = spliceTLV(0x36, strconv.Itoa(0))
	} else if status&0x10 == 0x10 {
		msg.ControlAction = spliceTLV(0x35, strconv.Itoa(0))
	}

	return msg.splice()
}

// 软件功能报文
func (msg *WHControlMessage) SoftFunction(option string) string {
	msg.ControlAction = spliceTLV(0x1d, option)
	return msg.splice()
}

// 热水器主控板特殊参数报文
func (msg *WHControlMessage) Special(option string) string {
	msg.ControlAction = spliceTLV(0x22, option)
	return msg.splice()
}
