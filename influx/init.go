package influx

// 待写入包
type Packet struct {
	SerialNumber      string
	CumulateHeatTime  int
	CumulateHotWater  int
	CumulateWorkTime  int
	CumulateUsedPower int
	CumulateSavePower int
}