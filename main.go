package main

import (
	"./mqtt"
	"./redis"
	"./protocol"
	"fmt"
)


func openMqtt() {
	m := new(mqtt.MQTT)

	m.Connect("zml-server", "tcp://192.168.2.108:1883")

	m.Subscribe("homeconsole", 2)
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
	protocol.Parse("Homeconsole02.010000018700140040000100100110238180717001012b000e100018071700010128000a0005000245")
}

func main() {
	fmt.Println("Start Point.")

	// testRedis()
	// openMqtt()
	testTlv()

	// for {}
}
