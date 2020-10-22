package receive

import "github.com/robertzml/Glaucus/tlv"

// 报文消息接口
// 所有接收的报文均实现改接口
type Message interface {
	// 报文协议解析
	// data: 返回的数据
	Parse(payload string) (data *tlv.TLV, err error)

	// 打印协议内容
	Print(cell tlv.TLV)

	// 安全检查
	Authorize(seq string) (pass bool)

	// 报文后续处理
	Handle(data *tlv.TLV, version float64, seq string) (err error)
}