package protocol

import "time"

// 校时报文
type TimingMessage struct {
	SerialNumber    string
	MainboardNumber string
	CurrentTime     string
}

// 拼接校时报文
func (msg *TimingMessage) splice() string {
	head := spliceHead()

	current := parseDateTimeToString(time.Now())

	sn := spliceTLV(0x127, msg.SerialNumber)
	mn := spliceTLV(0x12b, msg.MainboardNumber)
	st := spliceTLV(0x18, current)

	body := spliceTLV(0x11, sn + mn + st)

	return head + body
}