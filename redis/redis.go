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

// 读取数据
func (r *Redis) Read(key string) string {
	if val, err := redigo.String(r.Client.Do("GET", key)); err != nil {
		fmt.Println("read error: ", err.Error())
		return ""
	} else {
		return val
	}
}

// 检查key是否存在
// key: 键值
func (r *Redis) Exists(key string) int {
	result, err := r.Client.Do("EXISTS", key)
	if err != nil {
		fmt.Println("exist error: ", err.Error())
		return 0
	}

	return int(result.(int64))
}


 // 写入hash数据
 // key: 键值
 // s: 结构体
func (r *Redis) Hmset(key string, s interface{}) {
	if _, err := r.Client.Do("HMSET", redigo.Args{}.Add(key).AddFlat(s)...); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("redis update key:%s\n", key)
}

// 写入hash 中 某一项数据
func (r *Redis) Hset(key string, field string, val interface{}) {
	if _, err := r.Client.Do("HSET", key, field, val); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("redis update key:%s, field:%s, val: %v\n", key, field, val)
}

// 获取hash数据
// key: 键值
func (r *Redis) Hgetall(key string, dest interface{}) (err error) {
	v, err := redigo.Values(r.Client.Do("HGETALL", key))
	if err != nil {
		fmt.Printf("read equipment status failed.")
		return err
	}

	err = redigo.ScanStruct(v, dest);
	if err != nil {
		fmt.Printf("parse equipment status failed.")
		return err
	}

	return
}

/*
获取hash中一项的数据
 */
func (r *Redis) Hget(key string, field string) (result string) {
	reply, err := r.Client.Do("HGET", key, field)
	if (err != nil) {
		return ""
	}

	result = string(reply.([]byte))
	return
}
