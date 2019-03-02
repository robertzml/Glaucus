package equipment

const (
	// 实时状态前缀
	RealStatusPrefix = "real_"
)

type Equipment interface {
	// 从redis中获取设备状态
	GetStatus(serialnumber string) (err error)

	// 保持实时状态到redis中
	SaveStatus()
}

func SaveToRedis(equipment Equipment) {
	equipment.SaveStatus()
}