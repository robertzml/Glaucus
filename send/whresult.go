package send

import (
	"github.com/robertzml/Glaucus/tlv"
	"strconv"
)

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
	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)

	body := tlv.Splice(0x0011, sn + mn + msg.ResultAction)

	return head + body
}

// 设备重复
// D8 设备序列号重复
// D7 主板序列号重复
func (msg *WHResultMessage) Duplicate(option string) string {
	msg.ResultAction = tlv.Splice(0x13, option)
	return msg.splice()
}

// 快速响应
func (msg *WHResultMessage) Fast(option int) string {
	msg.ResultAction = tlv.Splice(0x16, strconv.FormatInt(int64(option), 16))
	return msg.splice()
}

// 设备响应周期
func (msg *WHResultMessage) Cycle(option int) string {
	msg.ResultAction = tlv.Splice(0x17, strconv.FormatInt(int64(option), 16))
	return msg.splice()
}

// 立即上报
func (msg *WHResultMessage) Reply() string {
	msg.ResultAction = tlv.Splice(0x19, "1")
	return msg.splice()
}