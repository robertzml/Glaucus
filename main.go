package main

import (
	"fmt"
	"./mqtt"
	"./redis"
	"./protocol"
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
	msg := "Homeconsole02.00000000010003006500010006123456001B00080000003F0007000110004002458783926-533E-484B-9B79-FEE11E5A6832001A00013001C00010"

	_, payload, err := protocol.ParseHead(msg)
	if err != nil {
		fmt.Println(err.Error())
	}

	tlv, err := protocol.ParseCell(payload)
	if err != nil {
		fmt.Println(err.Error())
	}

	if tlv.Tag != 0x03 {
		fmt.Println("tag incorrect")
	}

	fmt.Printf("Cell length: %d\n", tlv.Length)

	index := 0
	for index < tlv.Length {
		item, err := protocol.ParseTLV(tlv.Value, index);
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Tag: 0x%02X, Lenght: %d, Value: %s\n", item.Tag, item.Length, item.Value)
		index += item.Length + 8
	}
}

func main() {
	fmt.Println("Start Point.")

	// testRedis()
	// openMqtt()

	// for {}

	testTlv()
}
