package equipment

import (
	"encoding/json"
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

// 序列化数据
func serialize(v interface{}) string {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(data)
}