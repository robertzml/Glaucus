package main

import (
	"fmt"
	"time"
	"./mqtt"
	"./redis"
)


func openMqtt() {
	m := new(mqtt.MQTT)

	m.Connect("zml-server")

	m.Subscribe("earth")
	// m.Subscribe("world")
	m.Publish("earth", "this is another sample")


	time.Sleep(3 * time.Second)
}

func testRedis() {
	r := new(redis.Redis)
	r.Connect("192.168.2.116:6379")

	defer r.Close()

	r.Write("name", "bob")

	name := r.Read("name")
	fmt.Println(name)
}

func main() {
	fmt.Println("Start Point.")

	testRedis()
}
