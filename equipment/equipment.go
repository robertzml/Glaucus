package equipment

import (
	"encoding/json"
)

const (
	// 实时状态前缀
	RealStatusPrefix = "wh_"
)

type Equipment interface {
	// 从redis中获取设备状态
	LoadStatus(serialnumber string) (exists bool)

	// 保存实时状态到redis中
	SaveStatus()
}

func SaveToRedis(equipment Equipment) {
	equipment.SaveStatus()
}

// 序列化数据
func serialize(v interface{}) string {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(data)
}