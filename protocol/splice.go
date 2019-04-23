// 实现Homeconsole TLV 拼接工作

package protocol

import "fmt"

// 拼接设备控制报文
func (msg *WHControlMessage) splice() string {
	head := spliceHead()

	sn := spliceTLV(0x127, msg.SerialNumber)
	mn := spliceTLV(0x12b, msg.MainboardNumber)
	ca := spliceTLV(0x012, msg.ControlAction)

	body := spliceTLV(0x0010, sn+mn+ca)

	return head + body
}

// 拼接头部
func spliceHead() string {
	s := HomeConsoleVersion + fmt.Sprintf("%08x", 1)
	return s
}

// 拼接TLV
// tag: 信元编码
// val: 数据
// 返回：编码后的字符串
func spliceTLV(tag int, val string) string {
	return fmt.Sprintf("%04X%04X%s", tag, len(val), val)
}
