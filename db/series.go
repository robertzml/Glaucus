package db

// 时序数据存储接口
type Series interface {

	// 保存累积数据
	SaveCumulate(data interface{})
}