package equipment

import (
	"encoding/json"
)

const (
	// 热水器Redis前缀
	WaterHeaterPrefix = "wh_"

	// 净水器Redis前缀
	WaterCleanerPrefix = "wc_"
)

type Equipment interface {
	// 从redis中获取设备状态
	LoadStatus(serialNumber string) (exists bool)

	// 保存实时状态到redis中
	SaveStatus()
}


// 序列化数据
func serialize(v interface{}) string {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(data)
}