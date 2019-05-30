package equipment

import (
	"github.com/robertzml/Glaucus/redis"
)


// 热水器实时状态
type WaterHeater struct {
	SerialNumber      string
	MainboardNumber   string
	Logtime           int64
	DeviceType        string
	ControllerType    string
	Power             int8 // 开关状态
	OutTemp           int  // 出水温度
	OutFlow           int  // 出水流量
	ColdInTemp        int
	HotInTemp         int
	ErrorCode         int
	ErrorTime		  int64
	WifiVersion       string
	CumulateHeatTime  int
	CumulateHotWater  int
	CumulateWorkTime  int // 累计通电时间
	CumulateUsedPower int
	CumulateSavePower int
	Unlock            int8 // 解锁/加锁状态
	Activate          int8
	SetTemp           int
	SoftwareFunction  string
	OutputPower       int
	ManualClean       int8
	DeadlineTime      int64
	ActivationTime    int64
	SpecialParameter  string
	Online            int8
	LineTime          int64
	EnergySave		  int
	IMSI			  string
}

// 热水器设置状态
type WaterHeaterSetting struct {
	SerialNumber      	string
	SetActivateTime		int64
	Activate        	int8
	Unlock            	int8
	DeadlineTime    	int64
}

// 热水器运行数据
type WaterHeaterRunning struct {
	SerialNumber    string
	MainboardNumber string
	Logtime         int64
	Power           int8
	OutTemp         int
	OutFlow         int
	ColdInTemp      int
	HotInTemp       int
	SetTemp         int
	OutputPower     int
	ManualClean     int8
}

// 热水器报警数据
type WaterHeaterAlarm struct {
	SerialNumber    string
	MainboardNumber string
	Logtime         int64
	ErrorCode       int
	ErrorTime		int64
}

// 热水器关键数据
type WaterHeaterKey struct {
	SerialNumber    string
	MainboardNumber string
	Logtime         int64
	Activate        int8
	ActivationTime  int64
	Unlock          int8
	DeadlineTime    int64
	Online          int8
	LineTime        int64
}

// 热水器累计数据
type WaterHeaterCumulate struct {
	SerialNumber      string
	MainboardNumber   string
	Logtime           int64
	CumulateHeatTime  int
	CumulateHotWater  int
	CumulateWorkTime  int
	CumulateUsedPower int
	CumulateSavePower int
	ColdInTemp        int
	SetTemp           int
}

// 热水器登录数据
type WaterHeaterLogin struct {
	SerialNumber      string
	MainboardNumber   string
	Logtime           int64
	DeviceType        string
	ControllerType    string
	WifiVersion       string
	SoftwareFunction  string
}

// 获取redis中设备实时状态
// serialNumber: 设备序列号
// 返回 exists: 设备是否存在redis中
func (equipment *WaterHeater) LoadStatus(serialNumber string) (exists bool) {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	if rc.Exists(WaterHeaterPrefix + serialNumber) == 0 {
		return false
	}

	rc.Hgetall(WaterHeaterPrefix + serialNumber, equipment)

	return true
}

// 整体更新设备实时状态
func (equipment *WaterHeater) SaveStatus() {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	rc.Hmset(WaterHeaterPrefix+equipment.SerialNumber, equipment)
}

// 通过设备序列号获取主板序列号
func (equipment *WaterHeater) GetMainboardNumber(serialNumber string) (mainboardNumber string, exists bool) {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	mn := rc.Hget(WaterHeaterPrefix + serialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return "",false
	} else {
		return mn,true
	}
}

// 推送运行数据
func (equipment *WaterHeater) PushRunning(running *WaterHeaterRunning) {
	val := serialize(running)

	rc := new(redis.RedisClient)
	rc.Get(false)
	defer rc.Close()

	rc.Rpush(WaterHeaterPrefix + "running", val)
}

// 推送报警数据
func (equipment *WaterHeater) PushAlarm(alarm *WaterHeaterAlarm) {
	val := serialize(alarm)

	rc := new(redis.RedisClient)
	rc.Get(false)
	defer rc.Close()

	rc.Rpush(WaterHeaterPrefix + "alarm", val)
}

// 推送关键数据
func (equipment *WaterHeater) PushKey(key *WaterHeaterKey) {
	val := serialize(key)

	rc := new(redis.RedisClient)
	rc.Get(false)
	defer rc.Close()

	rc.Rpush(WaterHeaterPrefix + "key", val)
}

// 推送累计数据
func (equipment *WaterHeater) PushCumulate(cumulate *WaterHeaterCumulate) {
	val := serialize(cumulate)

	rc := new(redis.RedisClient)
	rc.Get(false)
	defer rc.Close()

	rc.Rpush(WaterHeaterPrefix + "cumulate", val)
}

// 推送登录数据
func (equipment *WaterHeater) PushLogin(login *WaterHeaterLogin) {
	val := serialize(login)

	rc := new(redis.RedisClient)
	rc.Get(false)
	defer rc.Close()

	rc.Rpush(WaterHeaterPrefix + "login", val)
}

// 获取设置状态
func (setting *WaterHeaterSetting) LoadSetting(serialNumber string) (exists bool) {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	if rc.Exists(WaterHeaterPrefix + "setting_" + serialNumber) == 0 {
		return false
	}

	rc.Hgetall(WaterHeaterPrefix + "setting_" + serialNumber, setting)

	return true
}

// 保存设置状态
func (setting *WaterHeaterSetting) SaveSetting() {
	rc := new(redis.RedisClient)
	rc.Get(true)
	defer rc.Close()

	rc.Hmset(WaterHeaterPrefix + "setting_" + setting.SerialNumber, setting)
}

