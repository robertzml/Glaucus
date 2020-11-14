package equipment

import (
	"encoding/json"
)

// 序列化数据
func serialize(v interface{}) string {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(data)
}
