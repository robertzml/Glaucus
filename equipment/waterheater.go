package equipment

import (
	"../redis"
	"errors"
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
func (equipment *WaterHeater) GetStatus(serialNumber string) (err error) {
	r := new(redis.Redis)
	defer r.Close()

	r.Connect()

	if r.Exists(RealStatusPrefix + serialNumber) == 0 {
		return errors.New("equipment not in redis.")
	}

	err = r.Hgetall(RealStatusPrefix + serialNumber, equipment)
	if err != nil {
		return
	}

	return nil
}

func (equipment *WaterHeater) SaveStatus() {

}