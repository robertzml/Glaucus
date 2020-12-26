package db

// 时序数据存储接口
type Series interface {

	// 保存累积数据
	SaveCumulate(tags map[string]string, fields map[string]interface{})

	// 保存基础数据
	SaveBasic(tags map[string]string, fields map[string]interface{})

	// 保存报警数据
	SaveAlarm(tags map[string]string, fields map[string]interface{})

	// 保存关键状态数据
	SaveKey(tags map[string]string, fields map[string]interface{})
}