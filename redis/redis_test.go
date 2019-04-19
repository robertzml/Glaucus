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

func TestPool(t *testing.T) {
	base.InitConfig()
	InitPool(0)

	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	v := rc.Hget("wh_01100101801100e2", "SerialNumber")
	fmt.Println(v)

	if r := recover(); r != nil {
		fmt.Printf("%v", r)
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