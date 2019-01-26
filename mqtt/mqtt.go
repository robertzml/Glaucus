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

// connect to mqtt server by clientId
func (m *MQTT) Connect(clientId string, address string) {
	m.ClientId = clientId
	m.Address = address

	opts := paho.NewClientOptions().AddBroker(address)
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(defaultHandler)

	m.client = paho.NewClient(opts)

	//create and start a client using the above ClientOptions
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// 关闭连接
func (m *MQTT) Disconnect() {
	m.client.Disconnect(250)
}

// 订阅相关主题，设置QoS
func (m *MQTT) Subscribe(topic string, qos byte, callback paho.MessageHandler) (err error) {
	if token := m.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		err = token.Error()
	} else {
		fmt.Printf("subscribe: %s\n", topic)
		err = nil
	}
	return
}

// 取消订阅
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
