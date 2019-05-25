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

// 拼接校时报文
func (msg *TimingMessage) Splice() string {
	head := tlv.SpliceHead()

	current := tlv.ParseDateTimeToString(time.Now())

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)
	st := tlv.Splice(0x18, current)

	body := tlv.Splice(0x11, sn + mn + st)

	return head + body
}