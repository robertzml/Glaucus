package redis

import "github.com/robertzml/Glaucus/equipment"

const (
	// 热水器Redis前缀
	waterHeaterPrefix = "wh_"
)

// redis 数据存储类
type Repository struct {
	client 	*RedisClient
}

// 获取redis中设备实时状态
// serialNumber: 设备序列号
// 返回 exists: 设备是否存在redis中
func (repo *Repository) LoadStatus(serialNumber string) (data *equipment.WaterHeater, exists bool) {
	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	if !rc.Exists(waterHeaterPrefix + serialNumber) {
		return nil, false
	}

	rc.Hgetall(waterHeaterPrefix + serialNumber, data)

	return data, true
}

// 保存热水器实时状态
func (repo *Repository) SaveStatus(data *equipment.WaterHeater) {
	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	rc.Hmset(waterHeaterPrefix+data.SerialNumber, data)
}

// 通过设备序列号获取主板序列号
func (repo *Repository) GetMainboardNumber(serialNumber string) (mainboardNumber string, exists bool) {
	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	mn := rc.Hget(waterHeaterPrefix + serialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return "",false
	} else {
		return mn,true
	}
}

// 设置 Redis {主板序列号 - 设备序列号} string
func SetMainboardString(mainboardNumber string, serialNumber string) {
	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	rc.Write(mainboardNumber, serialNumber)
}

// 读取 Redis {主板序列号 - 设备序列号} string
// 返回: 设备序列号
func (repo *Repository) GetMainboardString(mainboardNumber string) (serialNumber string) {
	rc := new(RedisClient)
	rc.Get()
	defer rc.Close()

	serialNumber, _ = rc.Read(mainboardNumber)
	return
}