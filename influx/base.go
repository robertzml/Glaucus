package influx

// 待写入包
type packet struct {
	SerialNumber		string
	MainboardNumber		string
	CumulateHeatTime	int
	CumulateHotWater	int
	CumulateWorkTime	int
	CumulateUsedPower	int
	CumulateSavePower	int
	EnergySave			int
}