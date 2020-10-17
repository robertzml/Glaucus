package base

// 接收到的MQTT 数据包
// 用于channel 同步
type ReceivePacket struct {
	ProductType int
	Topic 		string
	Payload 	string
}