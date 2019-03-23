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

func StartHttpServer(ch chan *base.SendPacket) {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         ":2450",
		WriteTimeout: 10 * time.Second,            //设置3秒的写超时
		Handler:      mux,
	}

	restHandler := new(RestHandler)
	restHandler.ch = ch

	mux.Handle("/", restHandler)
	mux.HandleFunc("/power", restHandler.power)
	mux.HandleFunc("/bye", sayBye)

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

		result := make(map[string]interface{})

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		serialNumber, ok := result["serialNumber"].(string)
		if !ok {
			w.WriteHeader(400)
			return
		}
		status, ok := result["status"].(float64)
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

			_, _ = io.WriteString(w, "control send")
		} else {
			_, _ = io.WriteString(w, "not found equipment")
		}

	} else {
		w.WriteHeader(404)
	}
}

func sayBye(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("bye bye ,this is v3 httpServer"))
	_, _ = io.WriteString(w, "say hi")
}