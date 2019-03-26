package rest

import (
	"../base"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTP接口处理结构体
type RestHandler struct {
	// 用于MQTT消息下发
	ch  chan *base.SendPacket
}

// HTTP返回消息
type ResponseMessage struct {
	Status		int		`json:"status"`
	Message		string	`json:"message"`
}

// 启动HTTP服务
func StartHttpServer(ch chan *base.SendPacket) {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         base.DefaultConfig.HttpListenAddress,
		WriteTimeout: 10 * time.Second,            //设置10秒的写超时
		Handler:      mux,
	}

	restHandler := new(RestHandler)
	restHandler.ch = ch

	mux.Handle("/", restHandler)
	mux.HandleFunc("/power", restHandler.power)
	mux.HandleFunc("/activate", restHandler.activate)
	mux.HandleFunc("/lock", restHandler.lock)
	mux.HandleFunc("/unlock", restHandler.unlock)
	mux.HandleFunc("/settemp", restHandler.setTemp)
	mux.HandleFunc("/deadline", restHandler.deadline)
	mux.HandleFunc("/clear", restHandler.clear)
	mux.HandleFunc("/clean", restHandler.clean)
	mux.HandleFunc("/special", restHandler.special)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("start server failed.")
	}
}

// 默认处理接口
func (*RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "hello")
}

// 返回消息
// status: 状态码
// message: 错误内容
func response(w http.ResponseWriter, status int, message string) {
	rm := ResponseMessage{ Status: status, Message: message }
	ret, _ := json.Marshal(&rm)
	_, _ = io.WriteString(w, string(ret))
}

