package rest

import (
	"encoding/json"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/equipment"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/protocol"
	"io/ioutil"
	"net/http"
	"time"
)

// 设备控制接口
func (handler *RestHandler) control(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func() {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			glog.Write(1, packageName, "control", err.Error())
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}

		typef, ok := parameter["type"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}
		controlType := int(typef)

		optionf, ok := parameter["option"].(float64)
		if !ok {
			response(w, 2, "option parameter is wrong.")
			return
		}
		option := int(optionf)

		control := new(protocol.WHControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber

			// 获取已保存的设置信息
			set := new(equipment.WaterHeaterSetting)
			_ = set.LoadSetting(serialNumber)
			set.SerialNumber = serialNumber

			switch controlType {
			case 1:
				pak.Payload = control.Power(option)
			case 2:
				pak.Payload = control.Activate(option)
				set.Activate = int8(option)
				if option == 1 {
					set.SetActivateTime = time.Now().Unix()
				}
			case 3:
				pak.Payload = control.Lock()
				set.Unlock = 0
			case 4:
				deadline, ok := parameter["deadline"].(float64)
				if !ok {
					w.WriteHeader(400)
					return
				}
				pak.Payload = control.Unlock(int64(deadline))

				set.Unlock = 1
				set.DeadlineTime = int64(deadline)
			case 5:
				deadline, ok := parameter["deadline"].(float64)
				if !ok {
					response(w, 3, "deadline parameter is wrong.")
					return
				}
				pak.Payload = control.SetDeadline(int64(deadline))
				set.DeadlineTime = int64(deadline)
			case 6:
				pak.Payload = control.SetTemp(option)
			case 7:
				pak.Payload = control.Clean(option)
			case 8:
				pak.Payload = control.Clean(option)
			default:
				w.WriteHeader(400)
				return
			}

			// 保存设置信息
			if controlType >=2 && controlType <= 5 {
				set.SaveSetting()
			}

			glog.Write(3, packageName, "control", "mqtt control producer.")
			base.MqttControlCh <- pak

			response(w, 0, "ok")

		} else {
			response(w, 1, "equipment not found.")
		}
	} else {
		w.WriteHeader(404)
	}
}

// 设备状态接口
func (handler *RestHandler) result(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func() {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			glog.Write(1, packageName, "result", err.Error())
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}

		typef, ok := parameter["type"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}
		resultType := int(typef)

		optionf, ok := parameter["option"].(float64)
		if !ok {
			response(w, 2, "option parameter is wrong.")
			return
		}
		option := int(optionf)

		waterHeader := new(equipment.WaterHeater)
		if mainboardNumber, exist := waterHeader.GetMainboardNumber(serialNumber); exist {
			resultMsg := protocol.NewWHResultMessage(serialNumber, mainboardNumber)

			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber

			switch resultType {
			case 1:
				pak.Payload = resultMsg.Fast(option)
			case 2:
				pak.Payload = resultMsg.Cycle(option)
			case 3:
				pak.Payload = resultMsg.Reply()
			default:
				w.WriteHeader(400)
				return
			}

			glog.Write(3, packageName, "control", "mqtt control producer.")
			base.MqttControlCh <- pak

			response(w, 0, "ok")

		} else {
			response(w, 1, "equipment not found.")
		}
	} else {
		w.WriteHeader(404)
	}
}

// 特殊参数设定
func (handler *RestHandler) special(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func() {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			glog.Write(1, packageName, "special", err.Error())
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}

		typef, ok := parameter["type"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}
		controlType := int(typef)

		option, ok := parameter["option"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.WHControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber

			switch controlType {
			case 1:
				pak.Payload = control.SoftFunction(option)
			case 2:
				pak.Payload = control.Special(option)
			default:
				w.WriteHeader(400)
				return
			}

			glog.Write(1, packageName, "special", "mqtt control producer.")
			base.MqttControlCh <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}
