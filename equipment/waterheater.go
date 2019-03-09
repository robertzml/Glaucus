package equipment

import (
	"../redis"
	"encoding/json"
	"fmt"
)

/*
热水器实时状态
 */
type WaterHeater struct {
	SerialNumber      string
	MainboardNumber   string
	Power             int8
	OutTemp           int
	OutFlow           int
	ColdInTemp        int
	HotInTemp         int
	ErrorCode         int
	WifiVersion       string
	CumulateHeatTime  int
	CumulateHotWater  int
	CumulateWorkTime  int
	CumulateUsedPower int
	CumualteSavePower int
	Lock              int8
	Activate          int8
	SetTemp           int
	SoftwareFunction  string
	OutputPower       int
	ManualClean       int8
	DeadlineTime      int64
	ActivationTime    int64
	SpecialParameter  string
}

// 获取redis中设备实时状态
// serialNumber: 设备序列号
// 返回 exists: 设备是否存在redis中
func (equipment *WaterHeater) GetStatus(serialNumber string) (exists bool, err error) {
	rc := new(redis.RedisClient)
	rc.Get()
	defer rc.Close()

	if rc.Exists(RealStatusPrefix + serialNumber) == 0 {
		return false, nil
	}

	err = rc.Hgetall(RealStatusPrefix + serialNumber, equipment)
	if err != nil {
		return true, err
	}

	return true,nil
}

// 整体更新设备实时状态
func (equipment *WaterHeater) SaveStatus() {
	rc := new(redis.RedisClient)
	rc.Get()
	defer rc.Close()

	rc.Hmset("real_" + equipment.SerialNumber, equipment)
}

// 部分更新设备实时状态
func (equipment *WaterHeater) UpdateField(field string, val interface{}) {
	rc := new(redis.RedisClient)
	rc.Get()
	defer rc.Close()

	rc.Hset(equipment.SerialNumber, field, val)
}

// 序列化设备属性
func (equipment *WaterHeater) Serialize() string {

	data, err := json.Marshal(equipment)
	if err != nil {
		fmt.Println(err)
	}

	return string(data)
}