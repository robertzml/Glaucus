package base

// 发送到MQTT数据包
// 用于channel 同步
type SendPacket struct {
	SerialNumber    string
	Payload			string
}
