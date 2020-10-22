package receive

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/equipment"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/send"
	"github.com/robertzml/Glaucus/tlv"
	"strconv"
	"time"
)

// 热水器设备状态报文
type WHStatusMessage struct {
	SerialNumber    string
	MainboardNumber string
	DeviceType      string
	ControllerType  string
}

// 解析协议内容
func (msg *WHStatusMessage) Parse(payload string) (data *tlv.TLV, err error) {
	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "whstatus parse", fmt.Sprintf("catch runtime panic: %v", r))
			err = fmt.Errorf("%v", r)
		}
	}()

	index := 0
	length := len(payload)

	for index < length {
		cell, err := tlv.ParseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus parse", fmt.Sprintf("error occur: %s", err.Error()))
			return nil, err
		}

		switch cell.Tag {
		case 0x127:
			msg.SerialNumber = cell.Value
		case 0x12b:
			msg.MainboardNumber = cell.Value
		case 0x125:
			msg.DeviceType = cell.Value
		case 0x12a:
			msg.ControllerType = cell.Value
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
func (msg *WHStatusMessage) Print(cell tlv.TLV) {
	fmt.Printf("Status Message Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

// 安全检查
// 返回: pass 是否通过
func (msg *WHStatusMessage) Authorize(seq string) (pass bool) {
	whs := new(equipment.WaterHeater)

	if exists := whs.LoadStatus(msg.SerialNumber); exists {
		if whs.MainboardNumber != msg.MainboardNumber {
			// 报文与redis缓存主板序列号不一致
			send.Write(msg.SerialNumber, msg.MainboardNumber, 1, "D8")

			glog.Write(2, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d8.", msg.SerialNumber, seq))
			return false
		}

		sn := equipment.GetMainboardString(whs.MainboardNumber)
		if len(sn) > 0 && sn != msg.SerialNumber {
			// 上报设备序列号与redis主板序列号-设备序列号映射 不一致
			send.Write(msg.SerialNumber, msg.MainboardNumber, 1, "D7")

			glog.Write(2, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d7.", msg.SerialNumber, seq))
			return false
		}

	} else { // 新设备
		sn := equipment.GetMainboardString(msg.MainboardNumber)
		if len(sn) > 0 && sn != msg.SerialNumber {
			// 主板序列号已存在
			send.Write(msg.SerialNumber, msg.MainboardNumber, 1, "D7")

			glog.Write(2, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d7 for new equipment.", msg.SerialNumber, seq))
			return false
		}

		glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. new equipment found.", msg.SerialNumber, seq))
		return true
	}

	glog.Write(4, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. pass.", msg.SerialNumber, seq))
	return true
}

// 报文后续处理
func (msg *WHStatusMessage) Handle(data *tlv.TLV, version float64, seq string) (err error) {
	var isFull bool
	if data.Tag == 0x128 {
		// 局部更新
		isFull = false
	} else if data.Tag == 0x12e {
		// 整体更新
		isFull = true
	} else {
		return errors.New("unknown tlv tag")
	}

	// 解析状态
	err, whs := msg.handleParseStatus(data.Value)
	if err != nil {
		return err
	}

	// 业务逻辑处理
	msg.handleLogic(whs, version, seq, isFull)

	glog.Write(5, packageName, "whstatus handle", fmt.Sprintf("sn: %s, seq: %s. handle finish.", msg.SerialNumber, seq))
	return nil
}

// 解析状态数据
// 返回：热水器状态
func (msg *WHStatusMessage) handleParseStatus(payload string) (err error, whs *equipment.WaterHeater) {
	whs = new(equipment.WaterHeater)

	index := 0
	length := len(payload)

	for index < length {
		cell, err := tlv.ParseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus parse status", fmt.Sprintf("sn: %s. error in parse tlv: %s", msg.SerialNumber, err.Error()))
			return err, nil
		}

		switch cell.Tag {
		case 0x01:
			v, _ := strconv.Atoi(cell.Value)
			whs.Power = int8(v)
		case 0x03:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.OutFlow = int(v)
		case 0x05:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.ErrorCode = int(v)
		case 0x08:
			whs.WifiVersion = cell.Value
		case 0x09:
			v, _ := tlv.ParseTime(cell.Value)
			whs.CumulateHeatTime = v
		case 0x0a:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulateHotWater = v
		case 0x0b:
			v, _ := tlv.ParseTime(cell.Value)
			whs.CumulateWorkTime = v
		case 0x0c:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulateUsedPower = v
		case 0x0d:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulateSavePower = v
		case 0x1a:
			v, _ := strconv.Atoi(cell.Value)
			whs.Unlock = int8(v)
		case 0x1b:
			v, _ := strconv.Atoi(cell.Value)
			whs.Activate = int8(v)
		case 0x1c:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.SetTemp = int(v)
		case 0x1d:
			whs.SoftwareFunction = cell.Value
		case 0x1e:
			v, _ := tlv.ParseCumulate(cell.Value, 4)
			whs.OutputPower = v
		case 0x1f:
			v, _ := strconv.Atoi(cell.Value)
			whs.ManualClean = int8(v)
		case 0x20:
			v, err := tlv.ParseDateToTimestamp(cell.Value)
			if err != nil {
				return err, nil
			}
			whs.DeadlineTime = v
		case 0x21:
			v, err := tlv.ParseDateToTimestamp(cell.Value)
			if err != nil {
				return err, nil
			}
			whs.ActivationTime = v
		case 0x22:
			whs.SpecialParameter = cell.Value
		case 0x23:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.EnergySave = int(v)
		case 0x24:
			whs.IMSI = cell.Value
		case 0x25:
			whs.ICCID = cell.Value
		case 0x26:
			whs.Coordinate = cell.Value
		case 0x27:
			whs.Csq = cell.Value
		}

		index += cell.Length + 8
	}

	return nil, whs
}

// 业务逻辑处理
// 参数： whs 解析出的新状态，保存whs 到 hash
func (msg *WHStatusMessage) handleLogic(whs *equipment.WaterHeater, version float64, seq string, isFull bool) {
	existsStatus := new(equipment.WaterHeater) // 原状态

	exists := existsStatus.LoadStatus(msg.SerialNumber)
	now := time.Now().Unix() * 1000

	// 全新设备 局部上报不处理
	if !exists && !isFull {
		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. cannot handle partial for new equipment.", msg.SerialNumber, seq))
		return
	}

	whs.SerialNumber = msg.SerialNumber
	whs.MainboardNumber = msg.MainboardNumber
	whs.Logtime = now
	whs.DeviceType = msg.DeviceType
	whs.ControllerType = msg.ControllerType
	whs.Online = 1

	// 整体上报
	if isFull {

		// 记录全上报时间
		whs.Fulltime = now

		glog.Write(3, packageName, "whstatus handle logic",
			fmt.Sprintf("sn: %s, seq: %s. full report, heat time: %d, hot water: %d, work time: %d, used power: %d, saved power: %d.",
				msg.SerialNumber, seq, whs.CumulateHeatTime, whs.CumulateHotWater, whs.CumulateWorkTime, whs.CumulateUsedPower, whs.CumulateSavePower))

		// influx.Write(whs.SerialNumber, whs.CumulateHeatTime, whs.CumulateHotWater, whs.CumulateWorkTime, whs.CumulateUsedPower, whs.CumulateSavePower)
	}

	// 全新设备整体上报
	if !exists && isFull {
		whs.LineTime = now

		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. new equipment find.", msg.SerialNumber, seq))

		whs.SaveStatus()
		return
	}

	// 后面开始处理已有设备
	whs.ErrorTime = existsStatus.ErrorTime
	whs.LineTime = existsStatus.LineTime

	// 更新 hash
	whs.SaveStatus()
}


