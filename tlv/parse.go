package tlv

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// 解析时间 转换为分钟
// 8位  最大FFFFFF小时+FF分钟
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


// 解析累积量
// 8位或4位  2位一转
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
		total = total * 100 + int(v)
	}

	return
}


// 解析日期 转换为 时间戳
// 10位 2位一转
// 输出13位
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

	timestamp = date.Unix() * 1000
	return
}

// 编码日期到TLV FFFFFFFFFF 格式
func ParseDateTimeToString(date time.Time) (string) {
	year := fmt.Sprintf("%02X", date.Year() - 2000)
	month := fmt.Sprintf("%02X", int(date.Month()))
	day := fmt.Sprintf("%02X", date.Day())
	hour := fmt.Sprintf("%02X", date.Hour())
	minute := fmt.Sprintf("%02X", date.Minute())

	return year + month + day + hour + minute
}

// 编码时间戳到TLV FFFFFFFFFF 格式
// timestamp 13位
func ParseTimestampToString(timestamp int64) (string) {
	timestamp = timestamp / 1000
	date := time.Unix(timestamp, 0)

	year := fmt.Sprintf("%02X", date.Year() - 2000)
	month := fmt.Sprintf("%02X", int(date.Month()))
	day := fmt.Sprintf("%02X", date.Day())
	hour := fmt.Sprintf("%02X", date.Hour())
	minute := fmt.Sprintf("%02X", date.Minute())

	return year + month + day + hour + minute
}
