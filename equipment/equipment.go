package equipment

import (
	"encoding/json"
)


// 设备数据库存储接口
type Context interface {

	// 连接数据库
	Connect()
}

// 设备实时信息存储接口
type Snapshot interface {

	// 当前状态
	Current()
}



// 序列化数据
func serialize(v interface{}) string {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(data)
}
