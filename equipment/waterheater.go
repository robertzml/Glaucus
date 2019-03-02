package equipment

/*
热水器上报状态
 */
type WaterHeater struct {
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
