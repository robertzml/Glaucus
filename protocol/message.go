package protocol


// 报文消息接口
// 所有类型的报文均实现改接口
type Message interface {
	// 报文协议解析
	// data: 返回的数据
	Parse(payload string) (data interface{}, err error)

	// 打印协议内容
	Print(cell TLV)

	// 安全检查
	Authorize() (pass bool, err error)

	// 报文后续处理
	Handle(data interface{}) (err error)
}