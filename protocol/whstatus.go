package protocol

import (
	"errors"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/equipment"
	"github.com/robertzml/Glaucus/glog"
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
		tlv, err := parseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus parse", fmt.Sprintf("error occur: %s", err.Error()))
			return nil, err
		}

		switch tlv.Tag {
		case 0x127:
			msg.SerialNumber = tlv.Value
		case 0x12b:
			msg.MainboardNumber = tlv.Value
		case 0x125:
			msg.DeviceType = tlv.Value
		case 0x12a:
			msg.ControllerType = tlv.Value
		default:
		}

		if tlv.Tag == 0x128 {
			return tlv, nil
		} else if tlv.Tag == 0x12e {
			return tlv, nil
		}

		index += tlv.Length + 8
	}

	return
}

// 打印协议信息
func (msg *WHStatusMessage) Print(cell TLV) {
	fmt.Printf("Status Message Print Tag: %#x, Serial Number:%s\n", cell.Tag, msg.SerialNumber)
}

// 安全检查
// 返回: pass 是否通过
func (msg *WHStatusMessage) Authorize() (pass bool, err error) {
	whs := new(equipment.WaterHeater)

	if exists := whs.LoadStatus(msg.SerialNumber); exists {
		if whs.MainboardNumber != msg.MainboardNumber {
			resMsg := NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.duplicate("D8")

			glog.Write(3, packageName, "whstatus authorize", "d8, mqtt control producer.")
			base.MqttControlCh <- pak

			return false, errors.New("mainboard number not equal.")
		}

		sn := equipment.GetMainboardString(whs.MainboardNumber)
		if (len(sn) > 0 && sn != msg.SerialNumber) {
			resMsg := NewWHResultMessage(msg.SerialNumber, msg.MainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = msg.SerialNumber
			pak.Payload = resMsg.duplicate("D7")

			glog.Write(3, packageName, "whstatus authorize", "d7, mqtt control producer.")
			base.MqttControlCh <- pak

			return false, errors.New("serial number not equal.")
		}
	} else {
		glog.Write(3, packageName, "whstatus authorize", "new equipment found.")
		return true, nil
	}

	glog.Write(3, packageName, "whstatus authorize", "pass.")
	return true, nil
}

// 报文后续处理
func (msg *WHStatusMessage) Handle(data interface{}) (err error) {
	switch data.(type) {
	case TLV:
		tlv := data.(TLV)
		if tlv.Tag == 0x128 {
			// 局部更新
			if err = msg.handleWaterHeaterChange(tlv.Value); err != nil {
				return err
			}
			glog.Write(3, packageName, "whstatus handle", "finish partial update.")
		} else if tlv.Tag == 0x12e {
			// 整体更新
			if err := msg.handleWaterHeaterTotal(tlv.Value); err != nil {
				return err
			}
			msg.timing()
			glog.Write(3, packageName, "whstatus handle", "finish total update.")
		}
	}

	if err := msg.handleSetting(); err != nil {
		return err
	}
	glog.Write(3, packageName, "whstatus handle", "setting compare pass.")

	return nil
}

// 整体解析热水器状态
func (msg *WHStatusMessage) handleWaterHeaterTotal(payload string) (err error) {
	waterHeaterStatus := new(equipment.WaterHeater)

	exists := waterHeaterStatus.LoadStatus(msg.SerialNumber)

	waterHeaterStatus.SerialNumber = msg.SerialNumber
	waterHeaterStatus.MainboardNumber = msg.MainboardNumber
	waterHeaterStatus.Logtime = time.Now().Unix() * 1000
	waterHeaterStatus.DeviceType = msg.DeviceType
	waterHeaterStatus.ControllerType = msg.ControllerType

	preErrorCode := waterHeaterStatus.ErrorCode
	preActivation := waterHeaterStatus.Activate

	index := 0
	length := len(payload)

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus handle", fmt.Sprintf("error in parse tlv: %s", err.Error()))
			return err
		}

		switch tlv.Tag {
		case 0x01:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.Power = int8(v)
		case 0x03:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.OutTemp = int(v)
		case 0x04:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.OutFlow = int(v)
		case 0x05:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.ColdInTemp = int(v)
		case 0x06:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.HotInTemp = int(v)
		case 0x07:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.ErrorCode = int(v)
		case 0x08:
			waterHeaterStatus.WifiVersion = tlv.Value
		case 0x09:
			v, _ := parseTime(tlv.Value)
			waterHeaterStatus.CumulateHeatTime = v
		case 0x0a:
			v, _ := parseCumulate(tlv.Value, 8)
			waterHeaterStatus.CumulateHotWater = v
		case 0x0b:
			v, _ := parseTime(tlv.Value)
			waterHeaterStatus.CumulateWorkTime = v
		case 0x0c:
			v, _ := parseCumulate(tlv.Value, 8)
			waterHeaterStatus.CumulateUsedPower = v
		case 0x0d:
			v, _ := parseCumulate(tlv.Value, 8)
			waterHeaterStatus.CumulateSavePower = v
		case 0x1a:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.Unlock = int8(v)
		case 0x1b:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.Activate = int8(v)
		case 0x1c:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			waterHeaterStatus.SetTemp = int(v)
		case 0x1d:
			waterHeaterStatus.SoftwareFunction = tlv.Value
		case 0x1e:
			v, _ := parseCumulate(tlv.Value, 4)
			waterHeaterStatus.OutputPower = v
		case 0x1f:
			v, _ := strconv.Atoi(tlv.Value)
			waterHeaterStatus.ManualClean = int8(v)
		case 0x20:
			v, _ := parseDateToTimestamp(tlv.Value)
			waterHeaterStatus.DeadlineTime = v
		case 0x21:
			v, _ := parseDateToTimestamp(tlv.Value)
			waterHeaterStatus.ActivationTime = v
		case 0x22:
			waterHeaterStatus.SpecialParameter = tlv.Value
		}

		index += tlv.Length + 8
	}

	if !exists || waterHeaterStatus.Online == 0 {
		waterHeaterStatus.LineTime = time.Now().Unix() * 1000

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

	if preErrorCode != waterHeaterStatus.ErrorCode {
		waterHeaterStatus.ErrorTime = time.Now().Unix() * 1000
	}

	if preActivation == 0 && waterHeaterStatus.Activate == 1 {
		msg.saveZeroCumulate()
	}

	waterHeaterStatus.Online = 1
	waterHeaterStatus.SaveStatus()

	equipment.SetMainboardString(waterHeaterStatus.MainboardNumber, waterHeaterStatus.SerialNumber)

	return
}

// 处理热水器变化状态，并局部更新
func (msg *WHStatusMessage) handleWaterHeaterChange(payload string) (err error) {
	whs := new(equipment.WaterHeater)

	exists := whs.LoadStatus(msg.SerialNumber)
	if !exists {
		glog.Write(2, packageName, "whstatus handle", "cannot update partial for new equipment.")
		return nil
	}

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

	cumulateChange := false

	index := 0
	length := len(payload)

	for index < length {
		tlv, err := parseTLV(payload, index)
		if err != nil {
			glog.Write(1, packageName, "whstatus handle", fmt.Sprintf("error in parse tlv: %s", err.Error()))
			return err
		}

		switch tlv.Tag {
		case 0x01:
			v, _ := strconv.Atoi(tlv.Value)
			whs.Power = int8(v)
			whRunning.Power = whs.Power
			runningChange = true
		case 0x03:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.OutTemp = int(v)
			whRunning.OutTemp = whs.OutTemp
			runningChange = true
		case 0x04:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.OutFlow = int(v)
			whRunning.OutFlow = whs.OutFlow
			runningChange = true
		case 0x05:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.ColdInTemp = int(v)
			whRunning.ColdInTemp = whs.ColdInTemp
			runningChange = true
			whCumulate.ColdInTemp = whs.ColdInTemp
			cumulateChange = true
		case 0x06:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.HotInTemp = int(v)
			whRunning.HotInTemp = whs.HotInTemp
			runningChange = true
		case 0x07:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.ErrorCode = int(v)

			if whAlarm.ErrorCode != whs.ErrorCode {
				whAlarm.ErrorTime = time.Now().Unix() * 1000
				whs.ErrorTime = whAlarm.ErrorTime
			}

			whAlarm.ErrorCode = whs.ErrorCode
			alarmChange = true
		case 0x08:
			whs.WifiVersion = tlv.Value
		case 0x09:
			v, _ := parseTime(tlv.Value)
			whs.CumulateHeatTime = v
			whCumulate.CumulateHeatTime = whs.CumulateHeatTime
			cumulateChange = true
		case 0x0a:
			v, _ := parseCumulate(tlv.Value, 8)
			whs.CumulateHotWater = v
			whCumulate.CumulateHotWater = whs.CumulateHotWater
			cumulateChange = true
		case 0x0b:
			v, _ := parseTime(tlv.Value)
			whs.CumulateWorkTime = v
			whCumulate.CumulateWorkTime = whs.CumulateWorkTime
			cumulateChange = true
		case 0x0c:
			v, _ := parseCumulate(tlv.Value, 8)
			whs.CumulateUsedPower = v
			whCumulate.CumulateUsedPower = whs.CumulateUsedPower
			cumulateChange = true
		case 0x0d:
			v, _ := parseCumulate(tlv.Value, 8)
			whs.CumulateSavePower = v
			whCumulate.CumulateSavePower = whs.CumulateSavePower
			cumulateChange = true
		case 0x1a:
			v, _ := strconv.Atoi(tlv.Value)
			whs.Unlock = int8(v)
			whKey.Unlock = whs.Unlock
			keyChange = true
		case 0x1b:
			v, _ := strconv.Atoi(tlv.Value)
			whs.Activate = int8(v)
			whKey.Activate = whs.Activate
			keyChange = true

			if preActivation == 0 && whs.Activate == 1 {
				msg.saveZeroCumulate()
			}
		case 0x1c:
			v, _ := strconv.ParseInt(tlv.Value, 16, 0)
			whs.SetTemp = int(v)
			whRunning.SetTemp = whs.SetTemp
			runningChange = true
			whCumulate.SetTemp = whs.SetTemp
			cumulateChange = true
		case 0x1d:
			whs.SoftwareFunction = tlv.Value
		case 0x1e:
			v, _ := parseCumulate(tlv.Value, 4)
			whs.OutputPower = v
			whRunning.OutputPower = whs.OutputPower
			runningChange = true
		case 0x1f:
			v, _ := strconv.Atoi(tlv.Value)
			whs.ManualClean = int8(v)
			whRunning.ManualClean = whs.ManualClean
			runningChange = true
		case 0x20:
			v, _ := parseDateToTimestamp(tlv.Value)
			whs.DeadlineTime = v
			whKey.DeadlineTime = whs.DeadlineTime
			keyChange = true
		case 0x21:
			v, _ := parseDateToTimestamp(tlv.Value)
			whs.ActivationTime = v
			whKey.ActivationTime = whs.ActivationTime
			keyChange = true
		case 0x22:
			whs.SpecialParameter = tlv.Value
		}

		index += tlv.Length + 8
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

// 处理比较设置数据
func (msg *WHStatusMessage) handleSetting() (err error) {
	whs := new(equipment.WaterHeater)

	exists := whs.LoadStatus(msg.SerialNumber)
	if !exists {
		glog.Write(2, packageName, "whstatus setting", "cannot compare setting for new equipment.")
		return nil
	}

	setting := new(equipment.WaterHeaterSetting)
	exists = setting.LoadSetting(msg.SerialNumber)
	if !exists {
		glog.Write(2, packageName, "whstatus setting", "setting is empty.")
		return nil
	}

	control := new(WHControlMessage)
	if ok := control.LoadEquipment(msg.SerialNumber); ok {
		pak := new(base.SendPacket)
		pak.SerialNumber = msg.SerialNumber

		if whs.Activate != setting.Activate {
			pak.Payload = control.Activate(int(setting.Activate))

			glog.Write(3, packageName, "whstatus setting", "activate, mqtt control producer.")
			base.MqttControlCh <- pak

			return nil
		}

		if setting.Activate == 0 {
			return nil
		}

		// 比较设备记录时间和设置激活时间，补发注销命令
		if whs.Activate == 1 && whs.ActivationTime+60*1000 < setting.SetActivateTime {
			pak.Payload = control.Activate(0)

			glog.Write(3, packageName, "whstatus setting", "inactivate, mqtt control producer.")
			base.MqttControlCh <- pak

			return nil
		}

		if whs.Unlock != setting.Unlock {
			if setting.Unlock == 0 {
				pak.Payload = control.Lock()

				glog.Write(3, packageName, "whstatus setting", "lock, mqtt control producer.")
				base.MqttControlCh <- pak
			} else {
				pak.Payload = control.Unlock(setting.DeadlineTime)

				glog.Write(3, packageName, "whstatus setting", "unlock, mqtt control producer.")
				base.MqttControlCh <- pak
			}

			return nil
		}

		if whs.DeadlineTime != setting.DeadlineTime {
			pak.Payload = control.SetDeadline(setting.DeadlineTime)

			glog.Write(3, packageName, "whstatus setting", "deadline, mqtt control producer.")
			base.MqttControlCh <- pak

			return nil
		}
	}

	return nil
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

// 下发校时
func (msg *WHStatusMessage) timing() {
	timing := new(TimingMessage)
	timing.SerialNumber = msg.SerialNumber
	timing.MainboardNumber = msg.MainboardNumber

	payload := timing.splice()

	pak := new(base.SendPacket)
	pak.SerialNumber = msg.SerialNumber
	pak.Payload = payload

	glog.Write(3, packageName, "whstatus timing", "send timing, mqtt control producer.")
	base.MqttControlCh <- pak
}
