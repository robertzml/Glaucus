package main

import (
	"./app"
	"fmt"
)

func main() {
	fmt.Println("Start Point.")

	app.StartMqtt()

	for {}
}
