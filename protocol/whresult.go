package protocol

import "strconv"

// 设备状态反馈报文
// 0x11
type WHResultMessage struct {
	SerialNumber    string
	MainboardNumber string
	ResultAction   string
}

// 设备状态反馈报文
func NewWHResultMessage(serialNumber string, mainboardNumber string) *WHResultMessage{
	return &WHResultMessage{ SerialNumber: serialNumber, MainboardNumber:mainboardNumber  }
}

// 拼接设备状态反馈报文
func (msg *WHResultMessage) splice() string {
	head := spliceHead()

	sn := spliceTLV(0x127, msg.SerialNumber)
	mn := spliceTLV(0x12b, msg.MainboardNumber)

	body := spliceTLV(0x0011, sn + mn + msg.ResultAction)

	return head + body
}

// 设备重复
// D8 设备序列号重复
// D7 主板序列号重复
func (msg *WHResultMessage) duplicate(option string) string {
	msg.ResultAction = spliceTLV(0x13, option)
	return msg.splice()
}

// 设备响应周期
func (msg *WHResultMessage) fast() string {
	msg.ResultAction = spliceTLV(0x16, "1")
	return msg.splice()
}

// 设备响应周期
func (msg *WHResultMessage) cycle(option int) string {
	msg.ResultAction = spliceTLV(0x17, strconv.FormatInt(int64(option), 16))
	return msg.splice()
}


// 立即上报
func (msg *WHResultMessage) reply() string {
	msg.ResultAction = spliceTLV(0x19, "1")
	return msg.splice()
}