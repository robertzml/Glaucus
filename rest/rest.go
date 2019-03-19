package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"../protocol"
)

func StartHttpServer() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         ":2450",
		WriteTimeout: 10 * time.Second,            //设置3秒的写超时
		Handler:      mux,
	}

	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/power", power)
	mux.HandleFunc("/bye", sayBye)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("start server failed.")
	}
}

type myHandler struct{}


func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "hello")
}


// 设备开关机接口
func power(w http.ResponseWriter, r *http.Request) {
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

		fmt.Println(status)
		control := new(protocol.ControlMessage)
		ok = control.LoadEquipment(serialNumber)
		fmt.Println(ok)

	} else {
		w.WriteHeader(404)
	}
}

func sayBye(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("bye bye ,this is v3 httpServer"))
	_, _ = io.WriteString(w, "say hi")
}