package redis

import (
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"testing"
)

/*
func TestSaveStruct(t *testing.T) {
	var w equipment.WaterHeater
	w.Power = 1
	w.OutTemp = 25
	w.OutFlow = 12

	r := new(Redis)
	r.Connect()

	r.Hmset("1234567", &w)

	r.Close()
}*/

//func TestPool(t *testing.T) {
//
//	base.InitConfig()
//	InitPool(0)
//
//	for i := 0; i < 3; i++ {
//
//		go func(num int) {
//			rc := new(RedisClient)
//			rc.Get()
//
//			fmt.Printf("time: %s, thread: %d, active: %d, idle: %d\n", time.Now(), num, RedisPool.ActiveCount(), RedisPool.IdleCount())
//
//			time.Sleep(10 * 1e9)
//			defer rc.Close()
//		}(i)
//
//
//		// fmt.Printf("time: %s, active: %d, idle: %d\n", time.Now(), RedisPool.ActiveCount(), RedisPool.IdleCount())
//	}
//
//	for i := 0; i < 5; i++ {
//		time.Sleep(10 * 1e9)
//		fmt.Printf("sleep, time: %s, active: %d, idle: %d\n", time.Now(), RedisPool.ActiveCount(), RedisPool.IdleCount())
//	}
//
//	r1 := new(RedisClient)
//	r1.Get()
//
//	fmt.Printf("time: %s, active: %d, idle: %d\n", time.Now(), RedisPool.ActiveCount(), RedisPool.IdleCount())
//	time.Sleep(10 * 1e9)
//
//	r1.Close()
//
//	fmt.Printf("time: %s, active: %d, idle: %d\n", time.Now(), RedisPool.ActiveCount(), RedisPool.IdleCount())
//
//	if r := recover(); r != nil {
//		fmt.Printf("%v", r)
//	}
//}


func TestDoGet(t *testing.T) {
	base.InitConfig()
	InitPool(0)

	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	if val, err := rc.Read("abc"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}


}

/*
func TestHset(t *testing.T) {
	rc := new(RedisClient)
	rc.Get()

	rc.Hset("real_01100101801100e2", "WifiVersion", "dse")


	rc.Close()
}*/