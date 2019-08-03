package receive

import (
	"fmt"
	"github.com/robertzml/Glaucus/equipment"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/tlv"
	"strconv"
	"time"
)

// 热水器离线消息报文
type WHOfflineMessage struct {
	SerialNumber    string
	MainboardNumber string
	Online			int
}


// 解析协议内容
func (msg *WHOfflineMessage) Parse(payload string) (data interface{}, err error) {
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
			return cell, nil
		} else if cell.Tag == 0x12e {
			return cell, nil
		}

		index += cell.Length + 8
	}

	return
}

// 打印协议信息
func (msg *WHOfflineMessage) Print(cell tlv.TLV) {
	fmt.Printf("Offline Message Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

// 安全检查
// 返回: pass 是否通过
func (msg *WHOfflineMessage) Authorize() (pass bool) {
	return true
}

// 报文后续处理
func (msg *WHOfflineMessage) Handle(data interface{}, seq string) (err error) {
	whs := new(equipment.WaterHeater)

	if exists := whs.LoadStatus(msg.SerialNumber); !exists {
		glog.Write(2, packageName, "whoffline handle", "cannot find equipment.")
		return nil
	}

	// 更新离线状态和时间
	whs.Online = 0
	whs.LineTime = time.Now().Unix() * 1000

	whs.SaveStatus()

	// 关键数据
	whKey := new(equipment.WaterHeaterKey)
	whKey.SerialNumber = whs.SerialNumber
	whKey.MainboardNumber = whs.MainboardNumber
	whKey.Logtime = whs.Logtime
	whKey.Activate = whs.Activate
	whKey.ActivationTime = whs.ActivationTime
	whKey.Unlock = whs.Unlock
	whKey.DeadlineTime = whs.DeadlineTime
	whKey.Online = 0
	whKey.LineTime = whs.LineTime

	whs.PushKey(whKey)

	glog.Write(3, packageName, "whoffline handle", fmt.Sprintf("sn: %s, seq: %s. save offline will message.", whs.SerialNumber, seq))
	return nil
}

