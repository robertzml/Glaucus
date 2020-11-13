package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/robertzml/Glaucus/base"
	"testing"
	"time"
)

func TestDoGet(t *testing.T) {
	base.LoadConfig()
	InitPool()

	rc := new(Client)
	rc.Open()
	defer rc.Close()

	rc.WriteString("abc", "hoop")

	if val, err := rc.ReadString("abc"); err != nil {
		t.Log(err)
	} else {
		t.Log(val)
	}
}

func newPool() *redigo.Pool {
	//return &redigo.Pool{
	//	MaxIdle: 3,
	//	IdleTimeout: 240 * time.Second,
	//	// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
	//	Dial: func () (redigo.Conn, error) { return redigo.Dial("tcp", "localhost",
	//		redigo.DialPassword("123"),
	//		) },
	//}

	// timeout := time.Duration(20)

	redisPool := &redigo.Pool{
		MaxIdle:         10,
		MaxActive:       50,
		IdleTimeout:     240 * time.Second,
		Wait:            true,
		MaxConnLifetime: 60 * time.Second,
		Dial: func() (redigo.Conn, error) {
			con, err := redigo.Dial("tcp", "localhost",
				redigo.DialPassword("123"))
			return con, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				fmt.Println(err)
			}
			return err
		},
	}

	return redisPool
}

func TestPool(t *testing.T) {
	pool := newPool()

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "jack", "tik tok")

	if err == nil {
		t.Log("this is ok")
	} else {
		t.Log("this is failed")
	}
}

/*
func TestHset(t *testing.T) {
	rc := new(RedisClient)
	rc.Get()

	rc.Hset("real_01100101801100e2", "WifiVersion", "dse")


	rc.Close()
}*/