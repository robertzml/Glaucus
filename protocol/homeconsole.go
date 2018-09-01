package protocol

import (
	"errors"
	"strconv"
)

const (
	HomeConsoleVersion = "Homeconsole02.00"
)

type Message interface {
	ParseContent(input string, pos int) TLV
}

/*
协议解析
根据收到的报文，解析出协议内容
 */
func Parse(message string) (err error){

	// read header
	_, payload, err := parseHead(message)
	if err != nil {
		return
	}

	cell, err := parseCell(payload)
	if err != nil {
		return
	}

	switch cell.Tag {
	case 0x03:
	default:
		err = errors.New("TLV not defined")
	}
	
	return
}

/*
解析协议头部
返回seq和协议内容
 */
func parseHead(message string) (seq string, payload string, err error) {
	vlen := len(HomeConsoleVersion)
	v := message[0: vlen]
	if HomeConsoleVersion != v {
		err = errors.New("version not match")
		return
	}

	seq = message[vlen: vlen + 8]
	payload = message[vlen + 8:]

	return
}

/*
解析信元
 */
func parseCell(payload string) (tlv TLV, err error) {
	tlv, err = parseTLV(payload, 0)
	return
}

/*
解析TLV
 */
func parseTLV(payload string, pos int) (tlv TLV, err error) {
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
