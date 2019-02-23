package equipment

/*
热水器上报状态
 */
type WaterHeater struct {
	Power		int
	OutTemp		int
	OutFlow		int
	ColdInTemp	int
	HotInTemp	int
	ErrorCode	int
}
