package mqtt

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Glaucus/glog"
)

// connect to mqtt server by clientId
func (m *MQTT) Connect(clientId string, username string, password string, address string, onConn paho.OnConnectHandler) {
	m.ClientId = clientId
	m.Address = address

	opts := paho.NewClientOptions().AddBroker(address)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(defaultHandler)
	opts.SetOnConnectHandler(onConn)

	m.client = paho.NewClient(opts)

	//create and start a client using the above ClientOptions
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		glog.WriteException(packageName, "Connect", token.Error().Error())
		panic(token.Error())
	}
}

// 关闭连接
func (m *MQTT) Disconnect() {
	m.client.Disconnect(250)
}

// 检查连接是否正常
func (m *MQTT) IsConnect() bool {
	if m.client == nil {
		return false
	} else {
		return m.client.IsConnected()
	}
}

// 订阅相关主题，设置QoS
func (m *MQTT) Subscribe(topic string, qos byte, callback paho.MessageHandler) (err error) {
	if token := m.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		err = token.Error()
	} else {
		glog.WriteInfo(packageName, "subscribe", "Topic:"+topic)
		err = nil
	}
	return
}

// 取消订阅
func (m *MQTT) Unsubscribe(topic string) (err error) {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		err = token.Error()
	}

	return err
}

// 发布订阅
func (m *MQTT) Publish(topic string, qos byte, payload string) {
	token := m.client.Publish(topic, qos, false, payload)
	token.Wait()
}
