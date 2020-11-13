package equipment

import (
	"github.com/robertzml/Glaucus/db"
)

const (
	// 热水器Redis前缀
	waterHeaterPrefix = "wh_"
)


// 热水器数据存储接口
type WaterHeaterRepo interface {

	// 保存热水器累计数据
	SaveCumulate(data *WaterHeaterCumulate)
}


// 热水器实时状态
type WaterHeater struct {
	SerialNumber      string
	MainboardNumber   string
	Logtime           int64
	Fulltime          int64 // 全上报时间
	DeviceType        string
	ControllerType    string
	Power             int8 // 开关状态
	OutTemp           int  // 出水温度
	OutFlow           int  // 出水流量
	ColdInTemp        int
	HotInTemp         int
	ErrorCode         int
	ErrorTime         int64 // 需设定
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
	Online            int8  // 需设定
	LineTime          int64 // 需设定
	EnergySave        int
	IMSI              string
	ICCID             string
	Coordinate        string
	Csq               string
}

// 热水器设置状态
type WaterHeaterSetting struct {
	SerialNumber    string
	SetActivateTime int64
	Activate        int8
	Unlock          int8
	DeadlineTime    int64
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
	ErrorTime       int64
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
	EnergySave        int
}

// 热水器登录数据
type WaterHeaterLogin struct {
	SerialNumber     string
	MainboardNumber  string
	Logtime          int64
	DeviceType       string
	ControllerType   string
	WifiVersion      string
	SoftwareFunction string
	ICCID            string
}

// 热水器数据异常
type WaterHeaterException struct {
	SerialNumber    string
	MainboardNumber string
	Logtime         int64
	Type            int
}

// 热水器数据处理类
type WaterHeaterHandler struct {
	// 实时数据操作接口
	snapshot db.Snapshot
}

func NewWaterHeaterHandler(snap db.Snapshot) *WaterHeaterHandler{
	handler := new(WaterHeaterHandler)
	handler.snapshot = snap

	return handler
}

// 获取redis中设备实时状态
// serialNumber: 设备序列号
// 返回 exists: 设备是否存在redis中
func (handler *WaterHeaterHandler) LoadStatus(serialNumber string) (data *WaterHeater, exists bool) {
	handler.snapshot.Open()
	defer handler.snapshot.Close()

	if !handler.snapshot.Exists(waterHeaterPrefix + serialNumber) {
		return nil, false
	}

	handler.snapshot.Load(waterHeaterPrefix + serialNumber, data)

	return data, true
}

// 保存热水器实时状态
func (handler *WaterHeaterHandler) SaveStatus(data *WaterHeater) {
	handler.snapshot.Open()
	defer handler.snapshot.Close()

	handler.snapshot.Save(waterHeaterPrefix+data.SerialNumber, data)
}

// 通过设备序列号获取主板序列号
func (handler *WaterHeaterHandler) GetMainboardNumber(serialNumber string) (mainboardNumber string, exists bool) {
	handler.snapshot.Open()
	defer handler.snapshot.Close()

	mn := handler.snapshot.LoadField(waterHeaterPrefix + serialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return "",false
	} else {
		return mn,true
	}
}

// 读取 Redis {主板序列号 - 设备序列号} string
// 返回: 设备序列号
func (handler *WaterHeaterHandler) GetMainboardString(mainboardNumber string) (serialNumber string) {
	handler.snapshot.Open()
	defer handler.snapshot.Close()

	serialNumber, _ = handler.snapshot.ReadString(mainboardNumber)
	return
}

// 设置 Redis {主板序列号 - 设备序列号} string
func (handler *WaterHeaterHandler) SetMainboardString(mainboardNumber string, serialNumber string) {
	handler.snapshot.Open()
	defer handler.snapshot.Close()

	handler.snapshot.WriteString(mainboardNumber, serialNumber)
}