package base

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	// 默认配置
	DefaultConfig Config

	// MQTT 发送控制指令 channel
	MqttControlCh  chan *SendPacket

	// MQTT 状态订阅消息 channel
	MqttStatusCh  chan *ReceivePacket
)

// 配置
type Config struct {
	// MQTT 服务器地址
	MqttServerAddress string

	// MQTT HTTP服务地址
	MqttServerHttp string

	// MQTT 用户名
	MqttUsername string

	// MQTT 密码
	MqttPassword string

	// MQTT主题订阅频道
	MqttChannel int

	// Rabbit MQ 连接字符串
	RabbitMQAddress string

	// Redis 数据库序号
	RedisDatabase int

	// Redis 服务器地址
	RedisServerAddress string

	// Redis 密码
	RedisPassword 	string

	// HTTP接口监听地址
	HttpListenAddress string

	// Influx地址
	InfluxAddress string

	// Influx Token
	InfluxToken string

	// Influx Org
	InfluxOrg string

	// 日志级别
	LogLevel	int

	// 输出日志到控制台
	LogToConsole bool
}

// 初始化默认配置
func InitConfig() {
	DefaultConfig.MqttServerAddress = "tcp://192.168.1.120:1883"
	DefaultConfig.MqttServerHttp = "http://192.168.1.120:18083"
	DefaultConfig.MqttChannel = 1
	DefaultConfig.MqttUsername = "glaucus"
	DefaultConfig.MqttPassword = "123456"
	DefaultConfig.RabbitMQAddress = "amqp://guest:guest@localhost:5672/"
	DefaultConfig.RedisDatabase = 0
	DefaultConfig.RedisServerAddress = "192.168.1.120:6379"
	DefaultConfig.RedisPassword = "123456"
	DefaultConfig.HttpListenAddress = ":8181"
	DefaultConfig.InfluxAddress = "127.0.0.1"
	DefaultConfig.InfluxToken = ""
	DefaultConfig.InfluxOrg = ""
	DefaultConfig.LogLevel = 3
	DefaultConfig.LogToConsole = true
}

// 载入配置文件
func LoadConfig()  {
	file, err := os.Open("./conf.json")
	if err != nil {
		fmt.Printf("cannot open the config file.\n")
		InitConfig()
		return
	}

	defer func() {
		_  = file.Close()
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&DefaultConfig)
	if err != nil {
		fmt.Printf("cannot parse the config file.\n")
		InitConfig()
		return
	}
}

// 初始化全局 channel
func InitChannel() {
	MqttControlCh = make(chan *SendPacket)
	MqttStatusCh = make(chan *ReceivePacket, 10)
}
