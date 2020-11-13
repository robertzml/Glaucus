package db

// 实时数据存储接口
type Snapshot interface {

	// 打开连接
	Open()

	// 关闭连接
	Close()

	// 写入字符串
	WriteString(key string, val string)

	// 读取字符串
	ReadString(key string) (string string, err error)

	// 检查key是否存在
	Exists(key string) bool

	// 保存对象数据
	Save(key string, s interface{})

	// 写入对象中某一项数据
	SaveField(key string, field string, val interface{})

	// 读取对象数据
	Load(key string, dest interface{})

	// 读取对象中一项的数据
	LoadField(key string, field string) (result string)
}