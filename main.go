package main

import (
	"fmt"
	"./mqtt"
	"time"
)

func main() {
	fmt.Println("Start Point.")

	m := new(mqtt.MQTT)

	m.Connect("zml-server")

	m.Subscribe("earth")
	// m.Subscribe("world")
	m.Publish("earth", "this is another sample")


	time.Sleep(3 * time.Second)

	for {
	}
}
