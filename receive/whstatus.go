package receive

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/base"
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
func (msg *WHStatusMessage) Parse(payload string) (data interface{}, err error) {
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
			return cell, nil
		} else if cell.Tag == 0x12e {
			return cell, nil
		}

		index += cell.Length + 8
	}

	return
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
		if whs.MainboardNumber != msg.MainboardNumber { // 主板序列号不一致
			resMsg := send.NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.Duplicate("D8")

			glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d8, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak

			return false
		}

		sn := equipment.GetMainboardString(whs.MainboardNumber)
		if len(sn) > 0 && sn != msg.SerialNumber { // 设备序列号不一致
			resMsg := send.NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.Duplicate("D7")

			glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d7, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak

			return false
		}
	} else {
		sn := equipment.GetMainboardString(msg.MainboardNumber)
		if len(sn) > 0 && sn != msg.SerialNumber { // 主板序列号已存在
			resMsg := send.NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.Duplicate("D7")

			glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. d7 for new equipment, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak

			return false
		}

		glog.Write(4, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. new equipment found.", msg.SerialNumber, seq))
		return true
	}

	glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s, seq: %s. pass.", msg.SerialNumber, seq))
	return true
}

// 报文后续处理
func (msg *WHStatusMessage) Handle(data interface{}, seq string) (err error) {
	switch data.(type) {
	case tlv.TLV:
		var isFull bool
		cell := data.(tlv.TLV)
		if cell.Tag == 0x128 {
			// 局部更新
			isFull = false
		} else if cell.Tag == 0x12e {
			// 整体更新
			isFull = true
		} else {
			return errors.New("unknown tlv tag")
		}

		// 解析状态
		err, whs := msg.handleParseStatus(cell.Value)
		if err != nil {
			return err
		}

		// 业务逻辑处理
		msg.handleLogic(whs, seq, isFull)

		// 比较设备设置状态
		if err := msg.handleSetting(seq); err != nil {
			return err
		}

		if isFull {
			// 设置 {主板序列号 - 设备序列号}
			equipment.SetMainboardString(msg.MainboardNumber, msg.SerialNumber)

			//校时
			msg.timing(seq)
		}

		glog.Write(3, packageName, "whstatus handle", fmt.Sprintf("sn: %s, seq: %s. handle finish.", msg.SerialNumber, seq))
		return nil

	default:
		// 无法进行后续处理
		return errors.New("wrong handle type.")
	}
}

// 解析状态数据
// 返回：热水器状态
func (msg* WHStatusMessage) handleParseStatus(payload string) (err error, whs *equipment.WaterHeater) {
	whs = new(equipment.WaterHeater)
	_ = whs.LoadStatus(msg.SerialNumber)

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
		}

		index += cell.Length + 8
	}

	return nil, whs
}

// 业务逻辑处理
// 参数： whs 解析出的新状态，保存whs 到 hash
func (msg* WHStatusMessage) handleLogic(whs *equipment.WaterHeater, seq string, isFull bool) {
	existsStatus := new(equipment.WaterHeater)	// 原状态

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

	// 全新设备整体上报
	if !exists && isFull {
		whs.LineTime = now

		if whs.ErrorCode != 0 {
			whs.ErrorTime = now

			glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. new equipment, push alarm.", msg.SerialNumber, seq))
			// 报警数据 推送 alarm list
			whAlarm := new(equipment.WaterHeaterAlarm)
			whAlarm.SerialNumber = whs.SerialNumber
			whAlarm.MainboardNumber = whs.MainboardNumber
			whAlarm.Logtime = whs.Logtime
			whAlarm.ErrorCode = whs.ErrorCode
			whAlarm.ErrorTime = whs.ErrorTime

			whs.PushAlarm(whAlarm)
		} else {
			whs.ErrorTime = 0
		}

		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. new equipment, push login.", msg.SerialNumber, seq))
		// 推送 login list
		whLogin := new(equipment.WaterHeaterLogin)
		whLogin.SerialNumber = whs.SerialNumber
		whLogin.MainboardNumber = whs.MainboardNumber
		whLogin.Logtime = now
		whLogin.DeviceType = whs.DeviceType
		whLogin.ControllerType = whs.ControllerType
		whLogin.WifiVersion = whs.WifiVersion
		whLogin.SoftwareFunction = whs.SoftwareFunction
		whLogin.ICCID = whs.ICCID

		whs.PushLogin(whLogin)

		whs.SaveStatus()

		return
	}


	// 后面开始处理已有设备
	whs.ErrorTime = existsStatus.ErrorTime
	whs.LineTime = existsStatus.LineTime

	// 设备重新上线，推送 wh_key list
	if existsStatus.Online == 0 {
		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. online, push key.", msg.SerialNumber, seq))

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

		whs.PushKey(whKey)
	}

	// 推送 wh_alarm list
	if whs.ErrorCode != 0 || existsStatus.ErrorCode != whs.ErrorCode {
		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push alarm.", msg.SerialNumber, seq))

		// 故障码变化，修改 ErrorTime
		if existsStatus.ErrorCode != whs.ErrorCode {
			whs.ErrorTime = now
		}

		whAlarm := new(equipment.WaterHeaterAlarm)
		whAlarm.SerialNumber = whs.SerialNumber
		whAlarm.MainboardNumber = whs.MainboardNumber
		whAlarm.Logtime = whs.Logtime
		whAlarm.ErrorCode = whs.ErrorCode
		whAlarm.ErrorTime = whs.ErrorTime

		whs.PushAlarm(whAlarm)
	}

	// 推送 running list
	if existsStatus.Power != whs.Power || existsStatus.OutTemp != whs.OutTemp || existsStatus.OutFlow != whs.OutFlow || existsStatus.ColdInTemp != whs.ColdInTemp ||
		existsStatus.HotInTemp != whs.HotInTemp || existsStatus.SetTemp != whs.SetTemp || existsStatus.OutputPower != whs.OutputPower ||
		existsStatus.ManualClean != whs.ManualClean {

		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push running.", msg.SerialNumber, seq))

		whRunning := new(equipment.WaterHeaterRunning)
		whRunning.SerialNumber = whs.SerialNumber
		whRunning.MainboardNumber = whs.MainboardNumber
		whRunning.Logtime = whs.Logtime
		whRunning.Power = whs.Power
		whRunning.OutTemp = whs.OutTemp
		whRunning.OutFlow = whs.OutFlow
		whRunning.ColdInTemp = whs.ColdInTemp
		whRunning.HotInTemp = whs.HotInTemp
		whRunning.SetTemp = whs.SetTemp
		whRunning.OutputPower = whs.OutputPower
		whRunning.ManualClean = whs.ManualClean

		whs.PushRunning(whRunning)
	}

	// 推送 key list
	if existsStatus.Unlock != whs.Unlock || existsStatus.Activate != whs.Activate || existsStatus.ActivationTime != whs.ActivationTime || existsStatus.DeadlineTime != whs.DeadlineTime {
		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push key.", msg.SerialNumber, seq))

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

		whs.PushKey(whKey)
	}

	// 推送 cumulate list
	if existsStatus.CumulateHeatTime != whs.CumulateHeatTime || existsStatus.CumulateHotWater != whs.CumulateHotWater || existsStatus.CumulateWorkTime != whs.CumulateWorkTime ||
		existsStatus.CumulateUsedPower != whs.CumulateUsedPower || existsStatus.CumulateSavePower != whs.CumulateSavePower || existsStatus.ColdInTemp != whs.ColdInTemp ||
		existsStatus.SetTemp != whs.SetTemp || existsStatus.EnergySave != whs.EnergySave {

		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push cumulate.", msg.SerialNumber, seq))

		whCumulate := new(equipment.WaterHeaterCumulate)
		whCumulate.SerialNumber = whs.SerialNumber
		whCumulate.MainboardNumber = whs.MainboardNumber
		whCumulate.Logtime = whs.Logtime
		whCumulate.CumulateHeatTime = whs.CumulateHeatTime
		whCumulate.CumulateHotWater = whs.CumulateHotWater
		whCumulate.CumulateWorkTime = whs.CumulateWorkTime
		whCumulate.CumulateUsedPower = whs.CumulateUsedPower
		whCumulate.CumulateSavePower = whs.CumulateSavePower
		whCumulate.ColdInTemp = whs.ColdInTemp
		whCumulate.SetTemp = whs.SetTemp
		whCumulate.EnergySave = whs.EnergySave

		whs.PushCumulate(whCumulate)
	}

	// 推送 login list
	if existsStatus.SoftwareFunction != whs.SoftwareFunction || existsStatus.WifiVersion != whs.WifiVersion || existsStatus.ICCID != whs.ICCID ||
		existsStatus.DeviceType != whs.DeviceType || existsStatus.ControllerType != whs.ControllerType {

		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push login.", msg.SerialNumber, seq))

		whLogin := new(equipment.WaterHeaterLogin)
		whLogin.SerialNumber = whs.SerialNumber
		whLogin.MainboardNumber = whs.MainboardNumber
		whLogin.Logtime = now
		whLogin.DeviceType = whs.DeviceType
		whLogin.ControllerType = whs.ControllerType
		whLogin.WifiVersion = whs.WifiVersion
		whLogin.SoftwareFunction = whs.SoftwareFunction
		whLogin.ICCID = whs.ICCID

		whs.PushLogin(whLogin)
	}

	// 检查数据异常
	if isFull && existsStatus.Activate == 1 && whs.Activate == 1 && whs.ActivationTime + 600 *1000 < now && (whs.CumulateHeatTime + 60 < existsStatus.CumulateHeatTime ||
		whs.CumulateHotWater + 120 < existsStatus.CumulateHotWater || whs.CumulateUsedPower + 200 < existsStatus.CumulateUsedPower ||
		whs.CumulateSavePower + 200 < existsStatus.CumulateSavePower) {

		glog.Write(3, packageName, "whstatus handle logic", fmt.Sprintf("sn: %s, seq: %s. push exception.", msg.SerialNumber, seq))

		whException := new(equipment.WaterHeaterException)
		whException.SerialNumber = whs.SerialNumber
		whException.MainboardNumber = whs.MainboardNumber
		whException.Logtime = whs.Logtime
		whException.Type = 1

		whs.PushException(whException)
	}


	// 已有设备从非激活态变为激活态，补零
	if existsStatus.Activate == 0 && whs.Activate == 1 {
		msg.saveZeroCumulate(seq)
	}

	// 更新 hash
	whs.SaveStatus()
}

// 补累计数据清零
func (msg *WHStatusMessage) saveZeroCumulate(seq string) {
	// 累计数据
	whCumulate := new(equipment.WaterHeaterCumulate)
	whCumulate.SerialNumber = msg.SerialNumber
	whCumulate.MainboardNumber = msg.MainboardNumber
	whCumulate.Logtime = time.Now().Unix() * 1000
	whCumulate.CumulateHeatTime = 0
	whCumulate.CumulateHotWater = 0
	whCumulate.CumulateWorkTime = 0
	whCumulate.CumulateUsedPower = 0
	whCumulate.CumulateSavePower = 0
	whCumulate.ColdInTemp = 0
	whCumulate.SetTemp = 0

	whs := new(equipment.WaterHeater)
	whs.PushCumulate(whCumulate)

	glog.Write(3, packageName, "whstatus handle", fmt.Sprintf("sn: %s, seq: %s, save zero cumulate.", msg.SerialNumber, seq))
}

// 处理比较设置数据
func (msg *WHStatusMessage) handleSetting(seq string) (err error) {
	whs := new(equipment.WaterHeater)

	exists := whs.LoadStatus(msg.SerialNumber)
	if !exists {
		glog.Write(2, packageName, "whstatus setting", fmt.Sprintf("sn: %s. cannot compare setting for new equipment.", msg.SerialNumber))
		return nil
	}

	setting := new(equipment.WaterHeaterSetting)
	exists = setting.LoadSetting(msg.SerialNumber)
	if !exists {
		glog.Write(2, packageName, "whstatus setting", fmt.Sprintf("sn: %s. setting is empty.", msg.SerialNumber))
		return nil
	}

	glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. set-active:%d, set-unlock:%d, set activate time: %d. status-active:%d, status-unlock:%d, status activate time: %d.",
		msg.SerialNumber, seq, setting.Activate, setting.Unlock, setting.SetActivateTime, whs.Activate, whs.Unlock, whs.ActivationTime))

	controlMsg := send.NewWHControlMessage(whs.SerialNumber, whs.MainboardNumber)

	pak := new(base.SendPacket)
	pak.SerialNumber = whs.SerialNumber

	if setting.Activate != whs.Activate {
		if setting.Activate == 0 { // whs.Activate == 1
			pak.Payload = controlMsg.Activate(0)

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. send inactivate, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak

		} else { // setting.Activate == 1 && whs.Activate == 0
			pak.Payload = controlMsg.Activate(1)

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. send activate, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak
		}

		return nil
	}

	// 比较设备记录时间和设置激活时间，补发注销命令
	if setting.Activate == 1 && whs.ActivationTime+360*1000 < setting.SetActivateTime {
		pak.Payload = controlMsg.Activate(0)

		glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. supply inactivate, MQTT control producer.", msg.SerialNumber, seq))
		base.MqttControlCh <- pak

		return nil
	}

	if setting.Activate == 0 && whs.Activate == 0 {
		return nil
	}

	if whs.Unlock != setting.Unlock {
		if setting.Unlock == 0 {
			pak.Payload = controlMsg.Lock()

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. lock, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak
		} else {
			pak.Payload = controlMsg.Unlock(1, setting.DeadlineTime)

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. unlock, MQTT control producer.", msg.SerialNumber, seq))
			base.MqttControlCh <- pak
		}

		return nil
	}

	if whs.DeadlineTime != setting.DeadlineTime {
		glog.Write(4, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. set-deadline:%d, status-deadline:%d.", msg.SerialNumber, seq, setting.DeadlineTime, whs.DeadlineTime))

		if setting.DeadlineTime == 0 {
			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. deadline is 0.", msg.SerialNumber, seq))
			return nil
		}

		pak.Payload = controlMsg.SetDeadline(setting.DeadlineTime)

		glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. deadline, MQTT control producer.", msg.SerialNumber, seq))
		base.MqttControlCh <- pak

		return nil
	}

	glog.Write(4, packageName, "whstatus setting", fmt.Sprintf("sn: %s, seq: %s. setting compare pass.", msg.SerialNumber, seq))
	return nil
}

// 下发校时
func (msg *WHStatusMessage) timing(seq string) {
	timing := send.NewTimingMessage(msg.SerialNumber, msg.MainboardNumber)
	payload := timing.Time()

	pak := new(base.SendPacket)
	pak.SerialNumber = msg.SerialNumber
	pak.Payload = payload

	glog.Write(4, packageName, "whstatus timing", fmt.Sprintf("sn: %s, seq: %s. send timing, MQTT control producer.", msg.SerialNumber, seq))
	base.MqttControlCh <- pak
}
