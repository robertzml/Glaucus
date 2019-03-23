package rest

import (
	"../base"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"../protocol"
)

type RestHandler struct {
	ch  chan *base.SendPacket
}

// HTTP返回消息
type ResponseMessage struct {
	Status		int
	Message		string
}

// 启动HTTP服务
func StartHttpServer(ch chan *base.SendPacket) {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         base.DefaultConfig.HttpListenAddress,
		WriteTimeout: 10 * time.Second,            //设置3秒的写超时
		Handler:      mux,
	}

	restHandler := new(RestHandler)
	restHandler.ch = ch

	mux.Handle("/", restHandler)
	mux.HandleFunc("/power", restHandler.power)
	mux.HandleFunc("/clear", restHandler.clear)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("start server failed.")
	}
}

func (*RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "hello")
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


// 返回消息
// status: 状态码
// message: 错误内容
func response(w http.ResponseWriter, status int, message string) {
	rm := ResponseMessage{ Status: status, Message: message }
	ret, _ := json.Marshal(&rm)
	_, _ = io.WriteString(w, string(ret))
}

