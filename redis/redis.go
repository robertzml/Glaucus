package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/robertzml/Glaucus/glog"
	"time"

	"github.com/robertzml/Glaucus/base"
)

const (
	packageName = "redis"
)

// redis 连接池
var RedisPool *redigo.Pool

// redis 连接
type RedisClient struct {
	client redigo.Conn
}

// 初始化Redis连接池
func InitPool() {
	timeout := time.Duration(20)

	RedisPool = &redigo.Pool{
		MaxIdle:         10,
		MaxActive:       50,
		IdleTimeout:     10 * time.Second,
		Wait:            true,
		MaxConnLifetime: 60 * time.Second,
		Dial: func() (redigo.Conn, error) {
			con, err := redigo.Dial("tcp", base.DefaultConfig.RedisServerAddress,
				redigo.DialPassword(base.DefaultConfig.RedisPassword),
				redigo.DialDatabase(base.DefaultConfig.RedisDatabase),
				redigo.DialConnectTimeout(timeout * time.Second),
				redigo.DialReadTimeout(timeout * time.Second),
				redigo.DialWriteTimeout(timeout * time.Second))
			if err != nil {
				// glog.Write(1, packageName, "dial", err.Error())
				fmt.Println("dial redis failed.")
				return nil, err
			}
			return con, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				glog.Write(1, packageName, "testOnBorrow", err.Error())
			}
			return err
		},
	}

	fmt.Println("redis pool create success.")
}

// 从连接池中获取一个redis 连接
func (r *RedisClient) Get() {
	r.client = RedisPool.Get()

	if r.client.Err() != nil {
		glog.Write(0, packageName, "get", r.client.Err().Error())
		panic(r.client.Err())
	}
	return
}

// 关闭连接
func (r *RedisClient) Close() {
	if err := r.client.Close(); err != nil {
		glog.Write(0, packageName, "close", err.Error())
		panic(err)
	}
}

// 写入数据
func (r *RedisClient) Write(key string, val string) {
	if _, err := r.client.Do("SET", key, val); err != nil {
		glog.Write(0, packageName, "write", err.Error())
		panic(err)
	}
}

// 读取数据
func (r *RedisClient) Read(key string) (string string, err error) {
	if val, err := redigo.String(r.client.Do("GET", key)); err != nil {
		return "", err
	} else {
		return val, nil
	}
}

// 检查key是否存在
// key: 键值
func (r *RedisClient) Exists(key string) bool {
	exists, err := redigo.Bool(r.client.Do("EXISTS", key))
	//result, err := r.client.Do("EXISTS", key)
	if err != nil {
		glog.Write(0, packageName, "exists", err.Error())
		panic(err)
	}

	// return int(result.(int64))
	return exists
}

// 写入hash数据
// key: 键值
// s: 结构体
func (r *RedisClient) Hmset(key string, s interface{}) {
	if _, err := r.client.Do("HMSET", redigo.Args{}.Add(key).AddFlat(s)...); err != nil {
		glog.Write(0, packageName, "hmset", err.Error())
		panic(err)
	}

	glog.Write(5, packageName, "hmset", fmt.Sprintf("redis update hash key:%s", key))
}

// 写入hash 中 某一项数据
func (r *RedisClient) Hset(key string, field string, val interface{}) {
	if _, err := r.client.Do("HSET", key, field, val); err != nil {
		glog.Write(0, packageName, "hset", err.Error())
		panic(err)
	}

	glog.Write(5, packageName, "hset", fmt.Sprintf("redis update key:%s, field:%s, val: %v", key, field, val))
}

// 获取hash数据
// key: 键值
// dest: 解析hash到指定结构体
func (r *RedisClient) Hgetall(key string, dest interface{}) {
	v, err := redigo.Values(r.client.Do("HGETALL", key))
	if err != nil {
		glog.Write(0, packageName, "hgetall", err.Error())
		panic(err)
	}

	if err = redigo.ScanStruct(v, dest); err != nil {
		glog.Write(0, packageName, "hgetall", err.Error())
		panic(err)
	}
}

// 获取hash中一项的数据
func (r *RedisClient) Hget(key string, field string) (result string) {
	reply, err := r.client.Do("HGET", key, field)
	if err != nil {
		glog.Write(0, packageName, "hget", err.Error())
		panic(err)
	}

	if reply == nil {
		return ""
	} else {
		result = string(reply.([]byte))
	}
	return
}

// 从右边推入队列
func (r *RedisClient) Rpush(key string, val string) {
	if _, err := r.client.Do("RPUSH", key, val); err != nil {
		glog.Write(0, packageName, "rpush", err.Error())
		panic(err)
	}
}
