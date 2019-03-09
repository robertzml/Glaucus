package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

const (
	RedisServer = "192.168.0.120:6379"
)

// redis 连接池
var RedisPool *redigo.Pool


// redis 连接
type RedisClient struct {
	client redigo.Conn
}

// 初始化Redis连接池
func InitPool() {
	timeout := time.Duration(30)

	// 建立连接池
	RedisPool = &redigo.Pool{
		MaxIdle:     100,
		MaxActive:   1000,
		IdleTimeout: 30 * time.Second,
		Wait:        true,
		Dial: func() (redigo.Conn, error) {
			con, err := redigo.Dial("tcp", RedisServer,
				redigo.DialConnectTimeout(timeout * time.Second),
				redigo.DialReadTimeout(timeout * time.Second),
				redigo.DialWriteTimeout(timeout * time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

// 从连接池中获取一个redis 连接
func (r *RedisClient) Get() {
	r.client = RedisPool.Get()
	if r.client.Err() != nil {
		panic(r.client.Err())
	} else {
		fmt.Println("redis connect allocate.")
	}
	return
}


// 关闭连接
func (r *RedisClient) Close() {
	err := r.client.Close()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("redis connect release.")
	}
}

// 写入数据
func (r *RedisClient) Write(key string, val string) bool {
	_, err := r.client.Do("SET", key, val)
	if err != nil {
		fmt.Println("write error. ", err.Error())
		return false
	} else {
		return true
	}
}

// 读取数据
func (r *RedisClient) Read(key string) string {
	if val, err := redigo.String(r.client.Do("GET", key)); err != nil {
		fmt.Println("read error: ", err.Error())
		return ""
	} else {
		return val
	}
}

// 检查key是否存在
// key: 键值
func (r *RedisClient) Exists(key string) int {
	result, err := r.client.Do("EXISTS", key)
	if err != nil {
		fmt.Println("exist error: ", err.Error())
		return 0
	}

	return int(result.(int64))
}


 // 写入hash数据
 // key: 键值
 // s: 结构体
func (r *RedisClient) Hmset(key string, s interface{}) {
	if _, err := r.client.Do("HMSET", redigo.Args{}.Add(key).AddFlat(s)...); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("redis update key:%s\n", key)
}

// 写入hash 中 某一项数据
func (r *RedisClient) Hset(key string, field string, val interface{}) {
	if _, err := r.client.Do("HSET", key, field, val); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("redis update key:%s, field:%s, val: %v\n", key, field, val)
}

// 获取hash数据
// key: 键值
func (r *RedisClient) Hgetall(key string, dest interface{}) (err error) {
	v, err := redigo.Values(r.client.Do("HGETALL", key))
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
func (r *RedisClient) Hget(key string, field string) (result string) {
	reply, err := r.client.Do("HGET", key, field)
	if (err != nil) {
		return ""
	}

	result = string(reply.([]byte))
	return
}
