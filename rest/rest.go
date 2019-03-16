package rest

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func StartHttpServer() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         ":2450",
		WriteTimeout: 10 * time.Second,            //设置3秒的写超时
		Handler:      mux,
	}

	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/bye", sayBye)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("start server failed.")
	}
}

type myHandler struct{}


func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "hello")
}


func sayBye(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("bye bye ,this is v3 httpServer"))
	_, _ = io.WriteString(w, "say hi")
}