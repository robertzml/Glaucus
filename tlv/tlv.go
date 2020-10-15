package tlv

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	HomeConsoleVersion = "Homeconsole04.00"
	packageName = "tlv"
)

type TLV struct {
	Tag    int
	Length int
	Value  string
}

func (tlv *TLV) String() string {
	return fmt.Sprintf("%04X%04X%s", tlv.Tag, tlv.Length, tlv.Value)
}

// TLV长度，包含头部
func (tlv *TLV) Size() int {
	return tlv.Length + 8
}


// 解析Homeconsole协议头部
// 返回seq和协议内容
func ParseHead(message string) (version float64, seq string, payload string, err error) {
	vlen := len(HomeConsoleVersion)
	v := message[0:vlen]
	version, _ = strconv.ParseFloat(v[11:], 32)

	if version < 4 {
		err = errors.New("version not match")
		return
	}

	seq = message[vlen : vlen+8]
	payload = message[vlen+8:]

	return
}

// 解析TLV
// payload: 需解析的字符串
// pos: 起始位置
func ParseTLV(payload string, pos int) (tlv TLV, err error) {
	tag, err := strconv.ParseInt(payload[pos:pos+4], 16, 0)
	if err != nil {
		return
	}
	tlv.Tag = int(tag)

	l, err := strconv.ParseInt(payload[pos+4:pos+8], 16, 0)
	if err != nil {
		return
	}
	tlv.Length = int(l)

	tlv.Value = payload[pos+8 : pos+8+tlv.Length]

	return
}

// 拼接Homeconsole头部
func SpliceHead() string {
	s := HomeConsoleVersion + fmt.Sprintf("%08x", 1)
	return s
}

// 拼接TLV
// tag: 信元编码
// val: 数据
// 返回：编码后的字符串
func Splice(tag int, val string) string {
	return fmt.Sprintf("%04X%04X%s", tag, len(val), val)
}
