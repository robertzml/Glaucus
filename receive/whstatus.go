package receive

import (
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
func (msg *WHStatusMessage) Authorize() (pass bool) {
	whs := new(equipment.WaterHeater)

	if exists := whs.LoadStatus(msg.SerialNumber); exists {
		if whs.MainboardNumber != msg.MainboardNumber { // 主板序列号不一致
			resMsg := send.NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.Duplicate("D8")

			glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s. d8, mqtt control producer.", msg.SerialNumber))
			base.MqttControlCh <- pak

			return false
		}

		sn := equipment.GetMainboardString(whs.MainboardNumber)
		if (len(sn) > 0 && sn != msg.SerialNumber) { // 设备序列号不一致
			resMsg := send.NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.Duplicate("D7")

			glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s. d7, mqtt control producer.", msg.SerialNumber))
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

			glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s. d7 for new equipment, mqtt control producer.", msg.SerialNumber))
			base.MqttControlCh <- pak

			return false
		}

		glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s. new equipment found.", msg.SerialNumber))
		return true
	}

	glog.Write(3, packageName, "whstatus authorize", fmt.Sprintf("sn: %s. pass.", msg.SerialNumber))
	return true
}

// 报文后续处理
func (msg *WHStatusMessage) Handle(data interface{}) (err error) {
	switch data.(type) {
	case tlv.TLV:
		cell := data.(tlv.TLV)
		if cell.Tag == 0x128 {
			// 局部更新
			if err = msg.handleWaterHeaterChange(cell.Value); err != nil {
				return err
			}
			glog.Write(3, packageName, "whstatus handle", fmt.Sprintf("sn: %s. finish partial update.", msg.SerialNumber))
		} else if cell.Tag == 0x12e {
			// 整体更新
			if err := msg.handleWaterHeaterTotal(cell.Value); err != nil {
				return err
			}
			msg.timing()
			glog.Write(3, packageName, "whstatus handle", fmt.Sprintf("sn: %s. finish total update.", msg.SerialNumber))
		}
	}

	// 比较设备设置状态
	if err := msg.handleSetting(); err != nil {
		return err
	}
	glog.Write(3, packageName, "whstatus handle", fmt.Sprintf("sn: %s. setting compare pass.", msg.SerialNumber))

	return nil
}

// 处理热水器变化状态，并局部更新
func (msg *WHStatusMessage) handleWaterHeaterChange(payload string) (err error) {
	whs := new(equipment.WaterHeater)

	exists := whs.LoadStatus(msg.SerialNumber)
	if !exists {
		glog.Write(2, packageName, "whstatus partial", fmt.Sprintf("sn: %s. cannot update partial for new equipment.", msg.SerialNumber))
		return nil
	}

	glog.Write(4, packageName, "whstatus partial", fmt.Sprintf("sn: %s. pre-active: %d, pre-online: %d, pre-unlock: %d.", msg.SerialNumber, whs.Activate, whs.Online, whs.Unlock))
	whs.Logtime = time.Now().Unix() * 1000

	preActivation := whs.Activate

	// 运行数据
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

	runningChange := false

	// 报警数据
	whAlarm := new(equipment.WaterHeaterAlarm)
	whAlarm.SerialNumber = whs.SerialNumber
	whAlarm.MainboardNumber = whs.MainboardNumber
	whAlarm.Logtime = whs.Logtime
	whAlarm.ErrorCode = whs.ErrorCode
	whAlarm.ErrorTime = whs.ErrorTime

	alarmChange := false

	// 关键数据
	keyChange := false

	whKey := new(equipment.WaterHeaterKey)
	whKey.SerialNumber = whs.SerialNumber
	whKey.MainboardNumber = whs.MainboardNumber
	whKey.Logtime = whs.Logtime
	whKey.Activate = whs.Activate
	whKey.ActivationTime = whs.ActivationTime
	whKey.Unlock = whs.Unlock
	whKey.DeadlineTime = whs.DeadlineTime

	if whs.Online == 0 {
		keyChange = true

		whs.LineTime = time.Now().Unix() * 1000
	}

	whs.Online = 1
	whKey.Online = whs.Online
	whKey.LineTime = whs.LineTime

	// 累计数据
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

	cumulateChange := false

	index := 0
	length := len(payload)

	for index < length {
		cell, err := tlv.ParseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus handle", fmt.Sprintf("error in parse tlv: %s", err.Error()))
			return err
		}

		switch cell.Tag {
		case 0x01:
			v, _ := strconv.Atoi(cell.Value)
			whs.Power = int8(v)
			whRunning.Power = whs.Power
			runningChange = true
		case 0x03:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.OutTemp = int(v)
			whRunning.OutTemp = whs.OutTemp
			runningChange = true
		case 0x04:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.OutFlow = int(v)
			whRunning.OutFlow = whs.OutFlow
			runningChange = true
		case 0x05:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.ColdInTemp = int(v)
			whRunning.ColdInTemp = whs.ColdInTemp
			runningChange = true
			whCumulate.ColdInTemp = whs.ColdInTemp
			cumulateChange = true
		case 0x06:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.HotInTemp = int(v)
			whRunning.HotInTemp = whs.HotInTemp
			runningChange = true
		case 0x07:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.ErrorCode = int(v)

			if whAlarm.ErrorCode != whs.ErrorCode {
				whAlarm.ErrorTime = time.Now().Unix() * 1000
				whs.ErrorTime = whAlarm.ErrorTime
			}

			whAlarm.ErrorCode = whs.ErrorCode
			alarmChange = true
		case 0x08:
			whs.WifiVersion = cell.Value
		case 0x09:
			v, _ := tlv.ParseTime(cell.Value)
			whs.CumulateHeatTime = v
			whCumulate.CumulateHeatTime = whs.CumulateHeatTime
			cumulateChange = true
		case 0x0a:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulateHotWater = v
			whCumulate.CumulateHotWater = whs.CumulateHotWater
			cumulateChange = true
		case 0x0b:
			v, _ := tlv.ParseTime(cell.Value)
			whs.CumulateWorkTime = v
			whCumulate.CumulateWorkTime = whs.CumulateWorkTime
			cumulateChange = true
		case 0x0c:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulateUsedPower = v
			whCumulate.CumulateUsedPower = whs.CumulateUsedPower
			cumulateChange = true
		case 0x0d:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			whs.CumulateSavePower = v
			whCumulate.CumulateSavePower = whs.CumulateSavePower
			cumulateChange = true
		case 0x1a:
			v, _ := strconv.Atoi(cell.Value)
			whs.Unlock = int8(v)
			whKey.Unlock = whs.Unlock
			keyChange = true
		case 0x1b:
			v, _ := strconv.Atoi(cell.Value)
			whs.Activate = int8(v)
			whKey.Activate = whs.Activate
			keyChange = true

			if preActivation == 0 && whs.Activate == 1 {
				msg.saveZeroCumulate()
			}
		case 0x1c:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.SetTemp = int(v)
			whRunning.SetTemp = whs.SetTemp
			runningChange = true
			whCumulate.SetTemp = whs.SetTemp
			cumulateChange = true
		case 0x1d:
			whs.SoftwareFunction = cell.Value
		case 0x1e:
			v, _ := tlv.ParseCumulate(cell.Value, 4)
			whs.OutputPower = v
			whRunning.OutputPower = whs.OutputPower
			runningChange = true
		case 0x1f:
			v, _ := strconv.Atoi(cell.Value)
			whs.ManualClean = int8(v)
			whRunning.ManualClean = whs.ManualClean
			runningChange = true
		case 0x20:
			v, _ := tlv.ParseDateToTimestamp(cell.Value)
			whs.DeadlineTime = v
			whKey.DeadlineTime = whs.DeadlineTime
			keyChange = true
		case 0x21:
			v, _ := tlv.ParseDateToTimestamp(cell.Value)
			whs.ActivationTime = v
			whKey.ActivationTime = whs.ActivationTime
			keyChange = true
		case 0x22:
			whs.SpecialParameter = cell.Value
		case 0x23:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			whs.EnergySave = int(v)
			whCumulate.EnergySave = whs.EnergySave
			cumulateChange = true
		case 0x24:
			whs.IMSI = cell.Value
		case 0x25:
			whs.ICCID = cell.Value
		}

		index += cell.Length + 8
	}

	whs.SaveStatus()

	if runningChange {
		whs.PushRunning(whRunning)
	}

	if alarmChange {
		whs.PushAlarm(whAlarm)
	}

	if keyChange {
		whs.PushKey(whKey)
	}

	if cumulateChange {
		whs.PushCumulate(whCumulate)
	}

	return nil
}

// 整体解析热水器状态
func (msg *WHStatusMessage) handleWaterHeaterTotal(payload string) (err error) {
	waterHeaterStatus := new(equipment.WaterHeater)

	exists := waterHeaterStatus.LoadStatus(msg.SerialNumber)
	now := time.Now().Unix() * 1000

	waterHeaterStatus.SerialNumber = msg.SerialNumber
	waterHeaterStatus.MainboardNumber = msg.MainboardNumber
	waterHeaterStatus.Logtime = now
	waterHeaterStatus.DeviceType = msg.DeviceType
	waterHeaterStatus.ControllerType = msg.ControllerType

	preErrorCode := waterHeaterStatus.ErrorCode
	preActivation := waterHeaterStatus.Activate
	cumulateChange := false

	index := 0
	length := len(payload)

	for index < length {
		cell, err := tlv.ParseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus total", fmt.Sprintf("sn: %s. error in parse tlv: %s", msg.SerialNumber, err.Error()))
			return err
		}

		switch cell.Tag {
		case 0x01:
			v, _ := strconv.Atoi(cell.Value)
			waterHeaterStatus.Power = int8(v)
		case 0x03:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			waterHeaterStatus.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			waterHeaterStatus.OutFlow = int(v)
		case 0x05:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			if int(v) != waterHeaterStatus.ColdInTemp {
				cumulateChange = true
			}
			waterHeaterStatus.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			waterHeaterStatus.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			waterHeaterStatus.ErrorCode = int(v)
		case 0x08:
			waterHeaterStatus.WifiVersion = cell.Value
		case 0x09:
			v, _ := tlv.ParseTime(cell.Value)
			if v != waterHeaterStatus.CumulateHeatTime {
				cumulateChange = true
			}
			waterHeaterStatus.CumulateHeatTime = v
		case 0x0a:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			if v != waterHeaterStatus.CumulateHotWater {
				cumulateChange = true
			}
			waterHeaterStatus.CumulateHotWater = v
		case 0x0b:
			v, _ := tlv.ParseTime(cell.Value)
			if v != waterHeaterStatus.CumulateWorkTime {
				cumulateChange = true
			}
			waterHeaterStatus.CumulateWorkTime = v
		case 0x0c:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			if v != waterHeaterStatus.CumulateUsedPower {
				cumulateChange = true
			}
			waterHeaterStatus.CumulateUsedPower = v
		case 0x0d:
			v, _ := tlv.ParseCumulate(cell.Value, 8)
			if v != waterHeaterStatus.CumulateSavePower {
				cumulateChange = true
			}
			waterHeaterStatus.CumulateSavePower = v
		case 0x1a:
			v, _ := strconv.Atoi(cell.Value)
			waterHeaterStatus.Unlock = int8(v)
		case 0x1b:
			v, _ := strconv.Atoi(cell.Value)
			waterHeaterStatus.Activate = int8(v)
		case 0x1c:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			if int(v) != waterHeaterStatus.SetTemp {
				cumulateChange = true
			}
			waterHeaterStatus.SetTemp = int(v)
		case 0x1d:
			waterHeaterStatus.SoftwareFunction = cell.Value
		case 0x1e:
			v, _ := tlv.ParseCumulate(cell.Value, 4)
			waterHeaterStatus.OutputPower = v
		case 0x1f:
			v, _ := strconv.Atoi(cell.Value)
			waterHeaterStatus.ManualClean = int8(v)
		case 0x20:
			v, _ := tlv.ParseDateToTimestamp(cell.Value)
			waterHeaterStatus.DeadlineTime = v
		case 0x21:
			v, _ := tlv.ParseDateToTimestamp(cell.Value)
			waterHeaterStatus.ActivationTime = v
		case 0x22:
			waterHeaterStatus.SpecialParameter = cell.Value
		case 0x23:
			v, _ := strconv.ParseInt(cell.Value, 16, 0)
			if int(v) != waterHeaterStatus.EnergySave {
				cumulateChange = true
			}
			waterHeaterStatus.EnergySave = int(v)
		case 0x24:
			waterHeaterStatus.IMSI = cell.Value
		case 0x25:
			waterHeaterStatus.ICCID = cell.Value
		}

		index += cell.Length + 8
	}

	// 全新设备，推送 wh_login list
	if !exists {
		glog.Write(3, packageName, "whstatus total", fmt.Sprintf("sn: %s. new equipment, push login", msg.SerialNumber))
		whLogin := new(equipment.WaterHeaterLogin)
		whLogin.SerialNumber = waterHeaterStatus.SerialNumber
		whLogin.MainboardNumber = waterHeaterStatus.MainboardNumber
		whLogin.Logtime = now
		whLogin.DeviceType = waterHeaterStatus.DeviceType
		whLogin.ControllerType = waterHeaterStatus.ControllerType
		whLogin.WifiVersion = waterHeaterStatus.WifiVersion
		whLogin.SoftwareFunction = waterHeaterStatus.SoftwareFunction

		waterHeaterStatus.PushLogin(whLogin)
	}

	// 在线状态变化，推送 wh_key list
	if !exists || waterHeaterStatus.Online == 0 {
		glog.Write(3, packageName, "whstatus total", fmt.Sprintf("sn: %s. online change. exists: %t, online: %d", msg.SerialNumber, exists, waterHeaterStatus.Online))

		waterHeaterStatus.LineTime = now

		whKey := new(equipment.WaterHeaterKey)
		whKey.SerialNumber = waterHeaterStatus.SerialNumber
		whKey.MainboardNumber = waterHeaterStatus.MainboardNumber
		whKey.Logtime = waterHeaterStatus.Logtime
		whKey.Activate = waterHeaterStatus.Activate
		whKey.ActivationTime = waterHeaterStatus.ActivationTime
		whKey.Unlock = waterHeaterStatus.Unlock
		whKey.DeadlineTime = waterHeaterStatus.DeadlineTime
		whKey.Online = 1
		whKey.LineTime = waterHeaterStatus.LineTime

		waterHeaterStatus.PushKey(whKey)
	}

	// 推送累计数据
	if cumulateChange {
		glog.Write(4, packageName, "whstatus total", fmt.Sprintf("sn: %s. push cumulative", msg.SerialNumber))

		whCumulate := new(equipment.WaterHeaterCumulate)
		whCumulate.SerialNumber = waterHeaterStatus.SerialNumber
		whCumulate.MainboardNumber = waterHeaterStatus.MainboardNumber
		whCumulate.Logtime = waterHeaterStatus.Logtime
		whCumulate.CumulateHeatTime = waterHeaterStatus.CumulateHeatTime
		whCumulate.CumulateHotWater = waterHeaterStatus.CumulateHotWater
		whCumulate.CumulateWorkTime = waterHeaterStatus.CumulateWorkTime
		whCumulate.CumulateUsedPower = waterHeaterStatus.CumulateUsedPower
		whCumulate.CumulateSavePower = waterHeaterStatus.CumulateSavePower
		whCumulate.ColdInTemp = waterHeaterStatus.ColdInTemp
		whCumulate.SetTemp = waterHeaterStatus.SetTemp
		whCumulate.EnergySave = waterHeaterStatus.EnergySave
		waterHeaterStatus.PushCumulate(whCumulate)
	}

	// 推送报警数据
	if preErrorCode != waterHeaterStatus.ErrorCode {
		whAlarm := new(equipment.WaterHeaterAlarm)
		whAlarm.SerialNumber = waterHeaterStatus.SerialNumber
		whAlarm.MainboardNumber = waterHeaterStatus.MainboardNumber
		whAlarm.Logtime = waterHeaterStatus.Logtime
		whAlarm.ErrorCode = waterHeaterStatus.ErrorCode
		whAlarm.ErrorTime = now

		waterHeaterStatus.PushAlarm(whAlarm)
	}

	// 已有设备从非激活态变为激活态，补零
	if exists && preActivation == 0 && waterHeaterStatus.Activate == 1 {
		msg.saveZeroCumulate()
	}

	waterHeaterStatus.Online = 1
	waterHeaterStatus.SaveStatus()

	equipment.SetMainboardString(waterHeaterStatus.MainboardNumber, waterHeaterStatus.SerialNumber)

	return
}

// 补累计数据清零
func (msg *WHStatusMessage) saveZeroCumulate() {
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

	glog.Write(3, packageName, "whstatus handle", "save zero cumulate.")
}

// 处理比较设置数据
func (msg *WHStatusMessage) handleSetting() (err error) {
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

	glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. set-active:%d, set-unlock:%d. status-active:%d, status-unlock:%d.",
		msg.SerialNumber, setting.Activate, setting.Unlock, whs.Activate, whs.Unlock))

	controlMsg := send.NewWHControlMessage(whs.SerialNumber, whs.MainboardNumber)

	pak := new(base.SendPacket)
	pak.SerialNumber = whs.SerialNumber

	if setting.Activate != whs.Activate {
		if setting.Activate == 0 { // whs.Activate == 1
			pak.Payload = controlMsg.Activate(0)

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. send inactivate, mqtt control producer.", msg.SerialNumber))
			base.MqttControlCh <- pak

		} else { // setting.Activate == 1 && whs.Activate == 0
			pak.Payload = controlMsg.Activate(1)

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. send activate, mqtt control producer.", msg.SerialNumber))
			base.MqttControlCh <- pak
		}

		return nil
	}

	// 比较设备记录时间和设置激活时间，补发注销命令
	if setting.Activate == 1 && whs.ActivationTime+360*1000 < setting.SetActivateTime {
		pak.Payload = controlMsg.Activate(0)

		glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. status activate time: %d, set activate time: %d",
			msg.SerialNumber, whs.ActivationTime, setting.SetActivateTime))
		glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. inactivate, mqtt control producer.", msg.SerialNumber))
		base.MqttControlCh <- pak

		return nil
	}

	if setting.Activate == 0 && whs.Activate == 0 {
		return nil
	}

	if whs.Unlock != setting.Unlock {
		if setting.Unlock == 0 {
			pak.Payload = controlMsg.Lock()

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. lock, mqtt control producer.", msg.SerialNumber))
			base.MqttControlCh <- pak
		} else {
			pak.Payload = controlMsg.Unlock(1, setting.DeadlineTime)

			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. unlock, mqtt control producer.", msg.SerialNumber))
			base.MqttControlCh <- pak
		}

		return nil
	}

	if whs.DeadlineTime != setting.DeadlineTime {
		glog.Write(4, packageName, "whstatus setting", fmt.Sprintf("sn: %s. set-deadline:%d, status-deadline:%d.", msg.SerialNumber, setting.DeadlineTime, whs.DeadlineTime))

		if setting.DeadlineTime == 0 {
			glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. deadline is 0.", msg.SerialNumber))
			return nil
		}

		pak.Payload = controlMsg.SetDeadline(setting.DeadlineTime)

		glog.Write(3, packageName, "whstatus setting", fmt.Sprintf("sn: %s. deadline, mqtt control producer.", msg.SerialNumber))
		base.MqttControlCh <- pak

		return nil
	}

	return nil
}

// 下发校时
func (msg *WHStatusMessage) timing() {
	timing := send.NewTimingMessage(msg.SerialNumber, msg.MainboardNumber)
	payload := timing.Time()

	pak := new(base.SendPacket)
	pak.SerialNumber = msg.SerialNumber
	pak.Payload = payload

	glog.Write(3, packageName, "whstatus timing", fmt.Sprintf("sn: %s. send timing, mqtt control producer.", msg.SerialNumber))
	base.MqttControlCh <- pak
}
