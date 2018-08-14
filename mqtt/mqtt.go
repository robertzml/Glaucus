package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	ClientId string
	client paho.Client
}


//define a function for the default message handler
var f paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// connect to mqtt server by clientId
func (m *MQTT) Connect(clientId string) {
	m.ClientId = clientId

	opts := paho.NewClientOptions().AddBroker("tcp://192.168.2.108:1883")
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(f)

	m.client = paho.NewClient(opts)

	//create and start a client using the above ClientOptions
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (m *MQTT) Disconnect() {
	m.client.Disconnect(250)
}

func (m *MQTT) Subscribe(topic string) {
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	} else {
		fmt.Printf("subscribe: %s\n", topic)
	}
}

func (m *MQTT) Unsubscribe(topic string) {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func (m *MQTT) Publish(topic string, payload string) {
	token := m.client.Publish(topic, 0, false, payload)
	token.Wait()
}

