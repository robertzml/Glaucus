package redis

import (
	redigo "github.com/gomodule/redigo/redis"
	"fmt"
)

type Redis struct {
	client redigo.Conn
}

// 连接服务器
func (r *Redis) Connect(server string) {
	rc, err := redigo.Dial("tcp", server)
	if err != nil {
		panic(err.Error())
	} else {
		r.client = rc
		fmt.Println("connect ok")
	}
}

// 关闭连接
func (r *Redis) Close() {
	r.client.Close()
}

// 写入数据
func (r *Redis) Write(key string, val string) bool {
	_, err := r.client.Do("SET", key, val)
	if err != nil {
		fmt.Println("write error. ", err.Error())
		return false
	} else {
		return true
	}
}

// 读取数据
func (r *Redis) Read(key string) string {
	if val, err := redigo.String(r.client.Do("GET", key)); err != nil {
		fmt.Println("read error. ", err.Error())
		return ""
	} else {
		return val
	}
}
