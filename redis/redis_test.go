package redis

import (
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

func TestPool(t *testing.T) {
	InitPool()

	r0 := new(RedisClient)
	r0.Get()
	r0.Write("abc", "sdfds")

	for i := 0; i < 100; i++ {
		rc := new(RedisClient)
		rc.Get()

		_ = rc.Read("abc")
		defer rc.Close()
	}
}

/*
func TestDoGet(t *testing.T) {
	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	val := rc.Read("abc")

	fmt.Println(val)
}*/

/*
func TestHset(t *testing.T) {
	rc := new(RedisClient)
	rc.Get()

	rc.Hset("real_01100101801100e2", "WifiVersion", "dse")


	rc.Close()
}*/