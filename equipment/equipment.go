package equipment

import (
	"encoding/json"
	"github.com/robertzml/Glaucus/redis"
)

const (
	// 热水器Redis前缀
	WaterHeaterPrefix = "wh_"

	// 净水器Redis前缀
	WaterCleanerPrefix = "wc_"
)

// 设备接口
type Equipment interface {
	// 从redis中获取设备状态
	LoadStatus(serialNumber string) (exists bool)

	// 保存实时状态到redis中
	SaveStatus()

	// 获取主板序列号
    GetMainboardNumber(serialNumber string) (mainboardNumber string, exists bool)
}

// 设置 Redis 主板序列号 string
func SetMainboardString(mainboardNumber string, serialNumber string) {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	rc.Write(mainboardNumber, serialNumber)
}

// 读取 Redis 主板序列号 string
// 返回: 设备序列号
func GetMainboardString(mainboardNumber string) (serialNumber string) {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	serialNumber, _ = rc.Read(mainboardNumber)
	return
}

// 序列化数据
func serialize(v interface{}) string {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(data)
}