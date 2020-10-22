package receive

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/tlv"
	"strconv"
)

// 热水器离线消息报文
type WHOfflineMessage struct {
	SerialNumber    string
	MainboardNumber string
	Online			int
}


// 解析协议内容
func (msg *WHOfflineMessage) Parse(payload string) (data *tlv.TLV, err error) {
	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "whoffline parse", fmt.Sprintf("catch runtime panic: %v", r))
			err = fmt.Errorf("%v", r)
		}
	}()

	index := 0
	length := len(payload)

	for index < length {
		cell, err := tlv.ParseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whoffline parse", fmt.Sprintf("error occur: %s", err.Error()))
			return nil, err
		}

		switch cell.Tag {
		case 0x127:
			msg.SerialNumber = cell.Value
		case 0x12b:
			msg.MainboardNumber = cell.Value
		case 0x129:
			msg.Online, _ = strconv.Atoi(cell.Value)
		default:
		}

		if cell.Tag == 0x128 {
			return &cell, nil
		} else if cell.Tag == 0x12e {
			return &cell, nil
		}

		index += cell.Length + 8
	}

	return nil, errors.New("cannot find info tag")
}

// 打印协议信息
func (msg *WHOfflineMessage) Print(cell tlv.TLV) {
	fmt.Printf("Offline Message Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

// 安全检查
// 返回: pass 是否通过
func (msg *WHOfflineMessage) Authorize(seq string) (pass bool) {
	return true
}

// 报文后续处理
func (msg *WHOfflineMessage) Handle(data *tlv.TLV, version float64, seq string) (err error) {
	glog.Write(3, packageName, "whoffline handle", fmt.Sprintf("sn: %s, seq: %s. save offline will message.", msg.SerialNumber, seq))
	return nil
}

