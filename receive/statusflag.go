package receive

type StatusChange int

const (
	IhPower				StatusChange = 1 << iota	// 1 << 0 is 0000000001
	IhOutTemp										// 1 << 1 is 0000000010
	IhOutFlow										// 1 << 2 is 0000000100
	IhColdInTemp									// 1 << 3 is 0000001000
	IhHotInTemp
	IhErrorCode
	IhWifiVersion
	IhCumulateHeatTime
	IhCumulateHotWater
	IhCumulateWorkTime
	IhCumulateUsedPower
	IhCumulateSavePower
	IhUnlock
	IhActivate
	IhSetTemp
	IhSoftwareFunction
	IhOutputPower
	IhManualClean
	IhDeadlineTime
	IhActivationTime
	IhSpecialParameter
	IhEnergySave
	IhIMSI
	IhICCID
)