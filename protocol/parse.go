// 实现Homeconsole TLV 解析工作，及内部数据格式转换。

package protocol

import (
	"errors"
	"strconv"
	"time"
)

/*
解析协议
根据收到的报文，解析出协议头部，确定协议类型
cell 报文头
msg  报文内容
 */
func parseType(message string) (cell TLV, msg Message, err error) {
	// read header
	_, payload, err := parseHead(message)
	if err != nil {
		return
	}

	// parse message cell type
	cell, err = parseTLV(payload, 0)
	if err != nil {
		return
	}

	switch cell.Tag {
	case 0x14:
		msg = new(StatusMessage)
	default:
		msg = nil
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
	v := message[0:vlen]
	if HomeConsoleVersion != v {
		err = errors.New("version not match")
		return
	}

	seq = message[vlen : vlen+8]
	payload = message[vlen+8:]

	return
}

/*
解析TLV
 */
func parseTLV(payload string, pos int) (tlv TLV, err error) {
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

/*
解析时间 转换为分钟
8位  最大FFFFFF小时+FF分钟
 */
func ParseTime(payload string) (totalMin int, err error) {
	if len(payload) != 8 {
		err = errors.New("time length is wrong.")
		return
	}

	hour, err := ParseCumulate(payload[0:6], 6)
	if err != nil {
		return
	}

	min, err := strconv.ParseInt(payload[6:8], 16, 0)
	if err != nil {
		return
	}

	totalMin = hour*60 + int(min)
	return
}

/*
解析累积量
8位或4位  2位一转
 */
func ParseCumulate(payload string, length int) (total int, err error) {
	if len(payload) != length {
		err = errors.New("cumulate length is wrong.")
		return
	}

	for i := 0; i < length; i += 2 {
		v, err := strconv.ParseInt(payload[i:i+2], 16, 0)
		if err != nil {
			break
		}
		total = total*100 + int(v)
	}

	return
}

/*
解析日期 转换为 时间戳
10位 2位一转
 */
func ParseDateToTimestamp(payload string) (timestamp int64, err error) {
	if len(payload) != 10 {
		err = errors.New("date length is wrong.")
		return
	}

	year, _ := strconv.ParseInt(payload[0:2], 16, 32)
	year += 2000

	month, _ := strconv.ParseInt(payload[2:4], 16, 0)
	day, _ := strconv.ParseInt(payload[4:6], 16, 0)
	hour, _ := strconv.ParseInt(payload[6:8], 16, 0)
	minute, _ := strconv.ParseInt(payload[8:10], 16, 0)

	date := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), 0, 0, time.Local)

	timestamp = date.Unix()
	return
}
