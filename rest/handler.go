package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"../base"
	"../protocol"
	"time"
)

func (handler *RestHandler) control(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
						

		} else {
			response(w, 1, "Equipment not found.")
		}
	} else {
		w.WriteHeader(404)
	}
}


// 设备开关机接口
func (handler *RestHandler) power(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		status, ok := parameter["status"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Power(int(status))

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 设备激活非激活
func (handler *RestHandler) activate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		status, ok := parameter["status"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Activate(int(status))

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 设备加锁
func (handler *RestHandler) lock(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]string)

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"]
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Lock()

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 设备解锁
func (handler *RestHandler) unlock(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		deadline, ok := parameter["deadline"].(time.Time)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Unlock(deadline)

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 设定温度
func (handler *RestHandler) setTemp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})
		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		temp, ok := parameter["temp"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.SetTemp(int(temp))

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 设置允许使用时间
func (handler *RestHandler) deadline(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		deadline, ok := parameter["deadline"].(time.Time)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.SetDeadline(deadline)

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 设备数据清零接口
// 参数status 使用位表示清零项目
func (handler *RestHandler) clear(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})
		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		status, ok := parameter["status"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Clear(int8(status))

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 手动清洗开关接口
func (handler *RestHandler) clean(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func () {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		status, ok := parameter["status"].(float64)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Clean(int(status))

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	} else {
		w.WriteHeader(404)
	}
}

// 热水器主控板特殊参数
func (handler *RestHandler) special(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		defer func() {
			_ = r.Body.Close()
		}()

		parameter := make(map[string]interface{})

		if err := json.Unmarshal(body, &parameter); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := parameter["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		option, ok := parameter["option"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}

		control := new(protocol.ControlMessage)
		if ok = control.LoadEquipment(serialNumber); ok {
			pak := new(base.SendPacket)
			pak.SerialNumber = serialNumber
			pak.Payload = control.Special(option)

			fmt.Println("control producer.")

			handler.ch <- pak

			response(w, 0, "ok")
		} else {
			response(w, 1, "Equipment not found.")
		}

	}else {
		w.WriteHeader(404)
	}
}
