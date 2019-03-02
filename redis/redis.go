package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	RedisServer = "192.168.0.120:6379"
)

type Redis struct {
	Client redigo.Conn
}

// 连接服务器
func (r *Redis) Connect() {
	rc, err := redigo.Dial("tcp", RedisServer)
	if err != nil {
		panic(err.Error())
	} else {
		r.Client = rc
		fmt.Println("redis connect ok")
	}
}

// 关闭连接
func (r *Redis) Close() {
	err := r.Client.Close()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("redis connect closed.")
	}
}

// 写入数据
func (r *Redis) Write(key string, val string) bool {
	_, err := r.Client.Do("SET", key, val)
	if err != nil {
		fmt.Println("write error. ", err.Error())
		return false
	} else {
		return true
	}
}

/*
整体写入设备实时状态
 */
func (r *Redis) Hmset(key string, s interface{}) {
	if _, err := r.Client.Do("HMSET", redigo.Args{}.Add(key).AddFlat(s)...); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("redis update key:%s\n", key)
}

// 读取数据
func (r *Redis) Read(key string) string {
	if val, err := redigo.String(r.Client.Do("GET", key)); err != nil {
		fmt.Println("read error. ", err.Error())
		return ""
	} else {
		return val
	}
}
