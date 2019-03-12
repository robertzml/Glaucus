package equipment

const (
	// 实时状态前缀
	RealStatusPrefix = "wh_"
)

type Equipment interface {
	// 从redis中获取设备状态
	GetStatus(serialnumber string) (exists bool, err error)

	// 保持实时状态到redis中
	SaveStatus()
}

func SaveToRedis(equipment Equipment) {
	equipment.SaveStatus()
}