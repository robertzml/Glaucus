package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	ClientId string
	Address  string
	client   paho.Client
}

//define a function for the default message handler
var f paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s, Id: %d, QoS: %d\n", msg.Topic(), msg.MessageID(), msg.Qos())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// connect to mqtt server by clientId
func (m *MQTT) Connect(clientId string, address string) {
	m.ClientId = clientId
	m.Address = address

	opts := paho.NewClientOptions().AddBroker(address)
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

func (m *MQTT) Subscribe(topic string, qos byte) (err error) {
	if token := m.client.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		err = token.Error()
	} else {
		fmt.Printf("subscribe: %s\n", topic)
		err = nil
	}
	return
}

func (m *MQTT) Unsubscribe(topic string) (err error){
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		err = token.Error()
	}

	return err
}

func (m *MQTT) Publish(topic string, qos byte, payload string) {
	token := m.client.Publish(topic, qos, false, payload)
	token.Wait()
}
