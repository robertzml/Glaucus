// 实现Homeconsole TLV 拼接工作

package protocol

import "fmt"


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
