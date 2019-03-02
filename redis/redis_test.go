package redis

import (
	"fmt"
	"testing"
	"../equipment"
)


func TestSaveStruct(t *testing.T) {
	var w equipment.WaterHeater
	w.Power = 1
	w.OutTemp = 25
	w.OutFlow = 12

	r := new(Redis)
	r.Connect()

	r.Hmset("1234567", &w)

	r.Close()
}

func TestDoGet(t *testing.T) {
	r := new(Redis)
	r.Connect()

	val := r.Read("abc")

	fmt.Println(val)
	r.Close()
}