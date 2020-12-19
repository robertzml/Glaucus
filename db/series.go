package db

// 时序数据存储接口
type Series interface {

	// 保存累积数据
	SaveCumulate(tags map[string]string, fields map[string]interface{})

	// 保存基础数据
	SaveBasic(tags map[string]string, fields map[string]interface{})
}