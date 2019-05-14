package base


var (
	// 默认配置
	DefaultConfig Config

	// MQTT 发送控制指令 channel
	MqttControlCh  chan *SendPacket

	// MQTT 接收上报状态 channel
	MqttStatusCh   chan *ReceivePacket
)

// 配置
type Config struct {
	// MQTT 服务器地址
	MqttServerAddress string

	// MQTT HTTP服务地址
	MqttServerHttp string

	// MQTT 用户名
	MqttUsername string

	// MQTT主题订阅频道
	MqttChannel int

	// Redis 数据库序号
	RedisDatabase int

	// Redis 服务器地址
	RedisServerAddress string

	// Redis 密码
	RedisPassword 	string

	// HTTP接口监听地址
	HttpListenAddress string
}

// 初始化默认配置
func InitConfig() {
	DefaultConfig.MqttServerAddress = "tcp://192.168.0.120:1883"
	DefaultConfig.MqttServerHttp = "http://192.168.0.120:18083"
	DefaultConfig.MqttChannel = 1
	DefaultConfig.MqttUsername = "glaucus"
	DefaultConfig.RedisDatabase = 0
	DefaultConfig.RedisServerAddress = "192.168.0.120:6379"
	DefaultConfig.RedisPassword = "123456"
	DefaultConfig.HttpListenAddress = ":8181"
}

// 初始化全局 channel
func InitChannel() {
	MqttControlCh = make(chan *SendPacket)
	MqttStatusCh = make(chan *ReceivePacket, 10)
}
