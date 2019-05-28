package rest

import (
	"encoding/json"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/equipment"
	"github.com/robertzml/Glaucus/glog"
	"github.com/robertzml/Glaucus/send"
	"io/ioutil"
	"net/http"
	"time"
)

// 设备控制请求参数
type ControlParam struct {
	SerialNumber string `json:"serialNumber"`
	Device 	     int	`json:"device"`
	ControlType	 int	`json:"type"`
	Option		 int	`json:"option"`
	Deadline	 int64	`json:"deadline"`
}

// 设备反馈请求参数
type ResultParam struct {
	SerialNumber string `json:"serialNumber"`
	Device 	     int	`json:"device"`
	ControlType	 int	`json:"type"`
	Option		 int	`json:"option"`
}

// 设备特殊控制请求参数
type SpecialParam struct {
	SerialNumber string `json:"serialNumber"`
	Device 	     int	`json:"device"`
	ControlType	 int	`json:"type"`
	Option		 string	`json:"option"`
}

// 设备控制接口
func (handler *RestHandler) control(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		body, _ := ioutil.ReadAll(r.Body)

		defer func() {
			_ = r.Body.Close()
		}()

		var param ControlParam
		if err := json.Unmarshal(body, &param); err != nil {
			w.WriteHeader(400)
			return
		}

		if param.Device == 1 { // 热水器
			status, msg, httpCode := handler.waterHeaterControl(param)

			if httpCode == 200 {
				response(w, status, msg)
			} else {
				w.WriteHeader(httpCode)
			}
		} else {
			response(w, 4, "unknown device.")
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

		var param ResultParam
		if err := json.Unmarshal(body, &param); err != nil {
			w.WriteHeader(400)
			return
		}

		if param.Device == 1 { // 热水器
			status, msg, httpCode := handler.waterHeaterResult(param)

			if httpCode == 200 {
				response(w, status, msg)
			} else {
				w.WriteHeader(httpCode)
			}
		} else {
			response(w, 4, "unknown device.")
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

		var param SpecialParam
		if err := json.Unmarshal(body, &param); err != nil {
			w.WriteHeader(400)
			return
		}

		if param.Device == 1 {
			waterHeater := new(equipment.WaterHeater)
			if mainboardNumber, exist := waterHeater.GetMainboardNumber(param.SerialNumber); exist {
				controlMsg := send.NewWHControlMessage(param.SerialNumber, mainboardNumber)

				pak := new(base.SendPacket)
				pak.SerialNumber = param.SerialNumber

				switch param.ControlType {
				case 1:
					pak.Payload = controlMsg.SoftFunction(param.Option)
				case 2:
					pak.Payload = controlMsg.Special(param.Option)
				default:
					w.WriteHeader(400)
					return
				}

				glog.Write(1, packageName, "special", "mqtt control producer.")
				base.MqttControlCh <- pak

				response(w, 0, "ok")

			} else {
				response(w, 1, "equipment not found.")
			}

		} else {
			response(w, 4, "unknown device.")
		}
	} else {
		w.WriteHeader(404)
	}
}

// 热水器控制内容
func (handler *RestHandler) waterHeaterControl(param ControlParam) (status int, msg string, httpCode int) {
	waterHeater := new(equipment.WaterHeater)

	if mainboardNumber, exist := waterHeater.GetMainboardNumber(param.SerialNumber); exist {
		controlMsg := send.NewWHControlMessage(param.SerialNumber, mainboardNumber)

		pak := new(base.SendPacket)
		pak.SerialNumber = param.SerialNumber

		// 获取已保存的设置信息
		set := new(equipment.WaterHeaterSetting)
		_ = set.LoadSetting(param.SerialNumber)
		set.SerialNumber = param.SerialNumber

		switch param.ControlType {
		case 1:
			pak.Payload = controlMsg.Power(param.Option)
		case 2:
			pak.Payload = controlMsg.Activate(param.Option)
			set.Activate = int8(param.Option)
			if param.Option == 1 {
				set.SetActivateTime = time.Now().Unix()
			}
		case 3:
			pak.Payload = controlMsg.Lock()
			set.Unlock = 0
		case 4:
			pak.Payload = controlMsg.Unlock(param.Deadline)
			set.Unlock = 1
			set.DeadlineTime = param.Deadline
		case 5:
			pak.Payload = controlMsg.SetDeadline(param.Deadline)
			set.DeadlineTime = param.Deadline
		case 6:
			pak.Payload = controlMsg.SetTemp(param.Option)
		case 7:
			pak.Payload = controlMsg.Clean(param.Option)
		case 8:
			pak.Payload = controlMsg.Clean(param.Option)
		default:
			return -1, "", 400
		}

		// 保存设置信息
		if param.ControlType >=2 && param.ControlType <= 5 {
			set.SaveSetting()
		}

		glog.Write(3, packageName, "control", "mqtt control producer.")
		base.MqttControlCh <- pak

		return 0, "ok", 200
	} else {
		return 1, "equipment not found.", 200
	}
}

// 热水器反馈内容
func (handler *RestHandler) waterHeaterResult(param ResultParam) (status int, msg string, httpCode int) {
	waterHeater := new(equipment.WaterHeater)
	if mainboardNumber, exist := waterHeater.GetMainboardNumber(param.SerialNumber); exist {
		resultMsg := send.NewWHResultMessage(param.SerialNumber, mainboardNumber)

		pak := new(base.SendPacket)
		pak.SerialNumber = param.SerialNumber

		switch param.ControlType {
		case 1:
			pak.Payload = resultMsg.Fast(param.Option)
		case 2:
			pak.Payload = resultMsg.Cycle(param.Option)
		case 3:
			pak.Payload = resultMsg.Reply()
		default:
			return -1, "", 400
		}

		glog.Write(3, packageName, "result", "mqtt control producer.")
		base.MqttControlCh <- pak

		return 0, "ok", 200
	} else {
		return 1, "equipment not found.", 200
	}
}