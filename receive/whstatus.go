package receive

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/db"
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

	Context *equipment.WaterHeaterContext
}

// 生成新热水器状态报文类
func NewWHStatusMessage(snapshot db.Snapshot, series db.Series) *WHStatusMessage {
	var msg = new(WHStatusMessage)

	msg.Context = equipment.NewWaterHeaterContext(snapshot, series)

	return msg
}

// 解析协议内容
func (msg *WHStatusMessage) Parse(payload string) (data *tlv.TLV, err error) {
	defer func() {
		if r := recover(); r != nil {
			glog.WriteError(packageName, "whstatus parse", fmt.Sprintf("catch runtime panic: %v", r))
			err = fmt.Errorf("%v", r)
		}
	}()

	index := 0
	length := len(payload)

	for index < length {
		cell, err := tlv.ParseTLV(payload, index)
		if err != nil {
			glog.WriteError(packageName, "whstatus parse", fmt.Sprintf("error occur: %s", err.Error()))
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

	if whs, exists := msg.Context.LoadStatus(msg.SerialNumber); exists {
		if whs.MainboardNumber != msg.MainboardNumber {
			// 报文与redis缓存主板序列号不一致
			send.WrteSpecial(msg.SerialNumber, 4, "D8")

			glog.WriteWarning(packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d8.", msg.SerialNumber, seq))
			return false
		}

		sn := msg.Context.GetMainboardString(whs.MainboardNumber)
		if len(sn) > 0 && sn != msg.SerialNumber {
			// 上报设备序列号与redis主板序列号-设备序列号映射 不一致
			send.WrteSpecial(msg.SerialNumber, 4, "D7")

			glog.WriteWarning(packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d7.", msg.SerialNumber, seq))
			return false
		}

	} else { // 新设备
		sn := msg.Context.GetMainboardString(msg.MainboardNumber)
		if len(sn) > 0 && sn != msg.SerialNumber {
			// 主板序列号已存在
			send.WrteSpecial(msg.SerialNumber, 4, "D7")

			glog.WriteWarning(packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d7 for new equipment.", msg.SerialNumber, seq))
			return false
		}

		glog.WriteInfo(packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. new equipment found.", msg.SerialNumber, seq))
		return true
	}

	glog.WriteDebug(packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. pass.", msg.SerialNumber, seq))
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

	if isFull {
		// 设置 {主板序列号 - 设备序列号}
		msg.Context.SetMainboardString(msg.MainboardNumber, msg.SerialNumber)

		//校时
		// msg.timing(seq)
	}

	glog.WriteVerbose(packageName, "whstatus handle", fmt.Sprintf("sn: %s, seq: %s. handle finish.", msg.SerialNumber, seq))
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
			glog.WriteError(packageName, "whstatus parse status", fmt.Sprintf("sn: %s. error in parse tlv: %s", msg.SerialNumber, err.Error()))
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
			whs.CumulativeHeatTime = v
		case 0x0a:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulativeHotWater = v
		case 0x0b:
			v, _ := tlv.ParseTime(cell.Value)
			whs.CumulativeWorkTime = v
		case 0x0c:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulativeUsedPower = v
		case 0x0d:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulativeSavePower = v
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

	existsStatus, exists := msg.Context.LoadStatus(msg.SerialNumber) // 原状态
	now := time.Now().Unix() * 1000

	// 全新设备 局部上报不处理
	if !exists && !isFull {
		glog.WriteDebug(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. cannot handle partial for new equipment.", msg.SerialNumber, seq))
		return
	}

	// 设置当前基础信息
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
	}

	// 全新设备整体上报
	if !exists && isFull {
		whs.LineTime = now

		glog.WriteDebug(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. new equipment find.", msg.SerialNumber, seq))

		// 处理错误状态
		if whs.ErrorCode != 0 {
			whs.ErrorTime = now

			glog.WriteWarning(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. new equipment, save alarm.", msg.SerialNumber, seq))

			// 报警数据 推送 alarm list
			whAlarm := new(equipment.WaterHeaterAlarm)
			whAlarm.SerialNumber = whs.SerialNumber
			whAlarm.MainboardNumber = whs.MainboardNumber
			whAlarm.Logtime = whs.Logtime
			whAlarm.ErrorCode = whs.ErrorCode
			whAlarm.ErrorTime = whs.ErrorTime

			msg.Context.SaveAlarm(whAlarm)
		} else {
			whs.ErrorTime = 0
		}

		glog.WriteInfo(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. new equipment, save basic and cumulate.", msg.SerialNumber, seq))

		// 保存 login list
		whBasic := new(equipment.WaterHeaterBasic)
		whBasic.SerialNumber = whs.SerialNumber
		whBasic.MainboardNumber = whs.MainboardNumber
		whBasic.Logtime = now
		whBasic.DeviceType = whs.DeviceType
		whBasic.ControllerType = whs.ControllerType
		whBasic.WifiVersion = whs.WifiVersion
		whBasic.SoftwareFunction = whs.SoftwareFunction
		whBasic.ICCID = whs.ICCID

		msg.Context.SaveBasic(whBasic)

		// 保存 cumulative list
		whCumulate := new(equipment.WaterHeaterCumulate)
		whCumulate.SerialNumber = whs.SerialNumber
		whCumulate.MainboardNumber = whs.MainboardNumber
		whCumulate.Logtime = now
		whCumulate.CumulativeHeatTime = whs.CumulativeHeatTime
		whCumulate.CumulativeHotWater = whs.CumulativeHotWater
		whCumulate.CumulativeWorkTime = whs.CumulativeWorkTime
		whCumulate.CumulativeUsedPower = whs.CumulativeUsedPower
		whCumulate.CumulativeSavePower = whs.CumulativeSavePower
		whCumulate.ColdInTemp = whs.ColdInTemp
		whCumulate.SetTemp = whs.SetTemp
		whCumulate.EnergySave = whs.EnergySave

		msg.Context.SaveCumulate(whCumulate)

		// 保存实时状态
		msg.Context.SaveStatus(whs)
		return
	}

	// 后面开始处理已有设备
	whs.ErrorTime = existsStatus.ErrorTime
	whs.LineTime = existsStatus.LineTime

	// 设备重新上线，推送 wh_key list
	if existsStatus.Online == 0 {
		glog.WriteInfo(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. online, save key status.", msg.SerialNumber, seq))

		whs.LineTime = now

		whKey := new(equipment.WaterHeaterKey)
		whKey.SerialNumber = whs.SerialNumber
		whKey.MainboardNumber = whs.MainboardNumber
		whKey.Logtime = whs.Logtime
		whKey.Activate = whs.Activate
		whKey.ActivationTime = whs.ActivationTime
		whKey.Unlock = whs.Unlock
		whKey.DeadlineTime = whs.DeadlineTime
		whKey.Online = 1
		whKey.LineTime = whs.LineTime

		msg.Context.SaveKey(whKey)
	}

	// 保存 wh_alarm
	if existsStatus.ErrorCode != whs.ErrorCode {
		glog.WriteWarning(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. save alarm.", msg.SerialNumber, seq))

		whs.ErrorTime = now

		whAlarm := new(equipment.WaterHeaterAlarm)
		whAlarm.SerialNumber = whs.SerialNumber
		whAlarm.MainboardNumber = whs.MainboardNumber
		whAlarm.Logtime = whs.Logtime
		whAlarm.ErrorCode = whs.ErrorCode
		whAlarm.ErrorTime = now

		msg.Context.SaveAlarm(whAlarm)
	}

	// 保存 key list
	if existsStatus.Unlock != whs.Unlock || existsStatus.Activate != whs.Activate || existsStatus.ActivationTime != whs.ActivationTime || existsStatus.DeadlineTime != whs.DeadlineTime {
		glog.WriteInfo(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. save key.", msg.SerialNumber, seq))

		whKey := new(equipment.WaterHeaterKey)
		whKey.SerialNumber = whs.SerialNumber
		whKey.MainboardNumber = whs.MainboardNumber
		whKey.Logtime = whs.Logtime
		whKey.Activate = whs.Activate
		whKey.ActivationTime = whs.ActivationTime
		whKey.Unlock = whs.Unlock
		whKey.DeadlineTime = whs.DeadlineTime
		whKey.Online = whs.Online
		whKey.LineTime = whs.LineTime

		msg.Context.SaveKey(whKey)
	}

	// 整体上报
	if isFull {
		// 保存 cumulative
		glog.WriteDebug(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push cumulate.", msg.SerialNumber, seq))

		whCumulate := new(equipment.WaterHeaterCumulate)
		whCumulate.SerialNumber = whs.SerialNumber
		whCumulate.MainboardNumber = whs.MainboardNumber
		whCumulate.Logtime = now
		whCumulate.CumulativeHeatTime = whs.CumulativeHeatTime
		whCumulate.CumulativeHotWater = whs.CumulativeHotWater
		whCumulate.CumulativeWorkTime = whs.CumulativeWorkTime
		whCumulate.CumulativeUsedPower = whs.CumulativeUsedPower
		whCumulate.CumulativeSavePower = whs.CumulativeSavePower
		whCumulate.ColdInTemp = whs.ColdInTemp
		whCumulate.SetTemp = whs.SetTemp
		whCumulate.EnergySave = whs.EnergySave

		msg.Context.SaveCumulate(whCumulate)

		// 保存 basic
		if existsStatus.SoftwareFunction != whs.SoftwareFunction || existsStatus.WifiVersion != whs.WifiVersion || existsStatus.ICCID != whs.ICCID ||
			existsStatus.DeviceType != whs.DeviceType || existsStatus.ControllerType != whs.ControllerType {

			glog.WriteDebug(packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. save basic.", msg.SerialNumber, seq))

			whBasic := new(equipment.WaterHeaterBasic)
			whBasic.SerialNumber = whs.SerialNumber
			whBasic.MainboardNumber = whs.MainboardNumber
			whBasic.Logtime = now
			whBasic.DeviceType = whs.DeviceType
			whBasic.ControllerType = whs.ControllerType
			whBasic.WifiVersion = whs.WifiVersion
			whBasic.SoftwareFunction = whs.SoftwareFunction
			whBasic.ICCID = whs.ICCID

			msg.Context.SaveBasic(whBasic)
		}
	}

	// 更新 hash
	msg.Context.SaveStatus(whs)
}
