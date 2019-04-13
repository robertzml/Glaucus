package equipment

import "github.com/robertzml/Glaucus/redis"

// 净水器实时状态
type WaterCleaner struct {
	SerialNumber      string
	MainboardNumber   string
	Logtime           int64
	DeviceType        string
	ControllerType    string

	Online            int8
	LineTime          int64
}

// 获取redis中净水器实时状态
// serialNumber: 设备序列号
// 返回 exists: 设备是否存在redis中
func (equipment *WaterCleaner) LoadStatus(serialNumber string) (exists bool) {
	rc := new(redis.RedisClient)
	rc.Get(1)
	defer rc.Close()

	if rc.Exists(WaterCleanerPrefix+serialNumber) == 0 {
		return false
	}

	rc.Hgetall(WaterCleanerPrefix+serialNumber, equipment)

	return true
}

// 整体更新设备实时状态
func (equipment *WaterCleaner) SaveStatus() {
	rc := new(redis.RedisClient)
	rc.Get(1)
	defer rc.Close()

	rc.Hmset(WaterCleanerPrefix+equipment.SerialNumber, equipment)
}