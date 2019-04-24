package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"time"

	"github.com/robertzml/Glaucus/base"
)

// redis 连接池
var RedisPool *redigo.Pool

// redis 连接
type RedisClient struct {
	client redigo.Conn
}


// 初始化Redis连接池
func InitPool(db int) {
	timeout := time.Duration(30)

	RedisPool = &redigo.Pool{
		MaxIdle:     100,
		MaxActive:   1000,
		IdleTimeout: 30 * time.Second,
		Wait:        true,
		Dial: func() (redigo.Conn, error) {
			con, err := redigo.Dial("tcp", base.DefaultConfig.RedisServerAddress,
				redigo.DialPassword(base.DefaultConfig.RedisPassword),
				redigo.DialDatabase(db),
				redigo.DialConnectTimeout(timeout*time.Second),
				redigo.DialReadTimeout(timeout*time.Second),
				redigo.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}



	// 建立连接池
	//RedisPool = &redigo.Pool{
	//	MaxIdle:     100,
	//	MaxActive:   1000,
	//	IdleTimeout: 30 * time.Second,
	//	Wait:        true,
	//	Dial: func() (redigo.Conn, error) {
	//		con, err := redigo.Dial("tcp", base.DefaultConfig.RedisServerAddress,
	//			redigo.DialPassword(base.DefaultConfig.RedisPassword),
	//			redigo.DialDatabase(1),
	//			redigo.DialConnectTimeout(timeout*time.Second),
	//			redigo.DialReadTimeout(timeout*time.Second),
	//			redigo.DialWriteTimeout(timeout*time.Second))
	//		if err != nil {
	//			return nil, err
	//		}
	//		return con, nil
	//	},
	//}
}

// 从连接池中获取一个redis 连接
// db: 数据库序号 0,1,2
func (r *RedisClient) Get() {
	r.client = RedisPool.Get()
	if r.client.Err() != nil {
		panic(r.client.Err())
	}
	return
}

// 关闭连接
func (r *RedisClient) Close() {
	if err := r.client.Close(); err != nil {
		panic(err.Error())
	}
}

// 写入数据
func (r *RedisClient) Write(key string, val string) {
	if _, err := r.client.Do("SET", key, val); err != nil {
		panic(err)
	}
}

// 读取数据
func (r *RedisClient) Read(key string) string {
	if val, err := redigo.String(r.client.Do("GET", key)); err != nil {
		panic(err)
	} else {
		return val
	}
}

// 检查key是否存在
// key: 键值
func (r *RedisClient) Exists(key string) int {
	result, err := r.client.Do("EXISTS", key)
	if err != nil {
		panic(err)
	}

	return int(result.(int64))
}

// 写入hash数据
// key: 键值
// s: 结构体
func (r *RedisClient) Hmset(key string, s interface{}) {
	if _, err := r.client.Do("HMSET", redigo.Args{}.Add(key).AddFlat(s)...); err != nil {
		panic(err)
	}
	fmt.Printf("redis update hash key:%s\n", key)
}

// 写入hash 中 某一项数据
func (r *RedisClient) Hset(key string, field string, val interface{}) {
	if _, err := r.client.Do("HSET", key, field, val); err != nil {
		panic(err)
	}
	fmt.Printf("redis update key:%s, field:%s, val: %v\n", key, field, val)
}

// 获取hash数据
// key: 键值
// dest: 解析hash到指定结构体
func (r *RedisClient) Hgetall(key string, dest interface{}) {
	v, err := redigo.Values(r.client.Do("HGETALL", key))
	if err != nil {
		panic(err)
	}

	if err = redigo.ScanStruct(v, dest); err != nil {
		panic(err)
	}
}

// 获取hash中一项的数据
func (r *RedisClient) Hget(key string, field string) (result string) {
	reply, err := r.client.Do("HGET", key, field)
	if err != nil {
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
		panic(err)
	}

	fmt.Printf("rpush key: %s\n", key)
}
