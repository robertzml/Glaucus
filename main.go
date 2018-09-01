package main

import (
	"./mqtt"
	"./redis"
	"fmt"
)


func openMqtt() {
	m := new(mqtt.MQTT)

	m.Connect("zml-server")

	m.Subscribe("homeconsole")
	// m.Subscribe("world")
	// m.Publish("earth", "this is another sample")


	// time.Sleep(3 * time.Second)
}

func testRedis() {
	r := new(redis.Redis)
	r.Connect("192.168.2.116:6379")

	defer r.Close()

	r.Write("name", "jim")

	name := r.Read("name")
	fmt.Println(name)
}

func testTlv() {

}

func main() {
	fmt.Println("Start Point.")

	// testRedis()
	// openMqtt()

	// for {}

	testTlv()
}
