package equipment

import (
	"fmt"
	"testing"
)

func TestReadWaterHeater(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v", r)
		}
	}()

	//base.LoadConfig()
	//redis.InitPool()

	//rc := new(redis.RedisClient)
	//rc.Get()
	//defer rc.Close()
	//
	//if !rc.Exists(waterHeaterPrefix + "01100101801100e2") {
	//	t.Log("this is false")
	//}
	//
	//whs := new(WaterHeater)
	//rc.Hgetall(waterHeaterPrefix + "01100101801100e2", whs)
	//
	//t.Log("this is true")

	//whs := new(WaterHeater)
	//
	//exists := whs.LoadStatus("01100101801100e2")

	//fmt.Printf("exists in redis: %v\n", exists)
	//fmt.Printf("%+v\n", whs)
}

func TestMainBoardNumber(t *testing.T) {
	//base.LoadConfig()
	//redis.InitPool()

	// sn := GetMainboardString("1111123")
	// t.Log(sn)
}