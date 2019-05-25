package send

// 报文消息接口
// 所有下发的报文均实现改接口
type Message interface {

	splice() string
}