package protocol

// 设备状态反馈报文
type WHResultMessage struct {
	SerialNumber    string
	MainboardNumber string

	// D8 设备序列号重复
	// D7 主板序列号重复
	ResultAction   string
}

// 拼接设备状态反馈报文
func (msg *WHResultMessage) splice() string {
	head := spliceHead()

	sn := spliceTLV(0x127, msg.SerialNumber)
	mn := spliceTLV(0x12b, msg.MainboardNumber)
	ra := spliceTLV(0x013, msg.ResultAction)

	body := spliceTLV(0x0011, sn + mn + ra)

	return head + body
}
