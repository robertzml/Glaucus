package base

// 默认配置
var DefaultConfig Config

// 配置
type Config struct {
	// MQTT 服务器地址
	MqttServerAddress string

	// MQTT HTTP服务地址
	MqttServerHttp string

	// MQTT主题订阅频道
	MqttChannel int

	// Redis 服务器地址
	RedisServerAddress string

	// HTTP接口监听地址
	HttpListenAddress string
}

// 初始化默认配置
func InitConfig() {
	DefaultConfig.MqttServerAddress = "tcp://192.168.0.120:1883"
	DefaultConfig.MqttServerHttp = "http://192.168.0.120:18083"
	DefaultConfig.MqttChannel = 1
	DefaultConfig.RedisServerAddress = "192.168.0.120:6379"
	DefaultConfig.HttpListenAddress = ":8181"
}
