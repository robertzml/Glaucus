package protocol

import (
	"errors"
	"fmt"
	"strconv"
)

var version string = "Homeconsole02.00"

type TLV struct {
	Tag    int
	Length int
	Value  string
}

func (tlv *TLV) String() string {
	return fmt.Sprintf("%04X%04X%s", tlv.Tag, tlv.Length, tlv.Value)
}


type Message interface {
	Parse(input string, pos int) TLV
}

/*
解析协议头部
返回seq和协议内容
 */
func ParseHead(message string) (seq string, payload string, err error) {
	vlen := len(version)
	v := message[0: vlen]
	if version != v {
		err = errors.New("version not match")
		return
	}

	seq = message[vlen: vlen + 8]
	payload = message[vlen + 8: len(message)]

	return
}

/*
解析信元
 */
func ParseCell(payload string) (tlv TLV, err error) {
	tlv, err = ParseTLV(payload, 0)
	return
}

/*
解析TLV
 */
func ParseTLV(payload string, pos int) (tlv TLV, err error) {
	tag, err := strconv.ParseInt(payload[pos: pos + 4], 16, 0)
	if err != nil {
		return
	}
	tlv.Tag = int(tag)

	l, err := strconv.ParseInt(payload[pos + 4: pos + 8], 16, 0)
	if err != nil {
		return
	}
	tlv.Length = int(l)

	tlv.Value = payload[pos + 8: pos + 8 + tlv.Length]

	return
}
