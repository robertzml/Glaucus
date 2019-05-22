package send

// 设备控制报文
// 0x10
type WHControlMessage struct {
	SerialNumber    string
	MainboardNumber string
	ControlAction   string
}

