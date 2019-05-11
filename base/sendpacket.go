package base

// 发送到MQTT数据包
// 用于channel 同步
type SendPacket struct {
	SerialNumber    string
	Payload			string
}

// 接收到的MQTT 数据包
// 用于channel 同步
type ReceivePacket struct {
	ProductType int
	Topic 		string
	Payload 	string
}