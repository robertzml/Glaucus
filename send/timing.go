package send

import (
	"github.com/robertzml/Glaucus/tlv"
	"time"
)

// 校时报文
type TimingMessage struct {
	SerialNumber    string
	MainboardNumber string
	CurrentTime     string
}

// 校时报文构造函数
func NewTimingMessage(serialNumber string, mainboardNumber string) *TimingMessage {
	return &TimingMessage{SerialNumber: serialNumber, MainboardNumber: mainboardNumber}
}

// 拼接校时报文
func (msg *TimingMessage) splice() string {
	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)
	ct := tlv.Splice(0x18, msg.CurrentTime)

	body := tlv.Splice(0x11, sn + mn + ct)

	return head + body
}

// 校时
func (msg *TimingMessage) Time() string {
	msg.CurrentTime = tlv.ParseDateTimeToString(time.Now())
	return msg.splice()
}