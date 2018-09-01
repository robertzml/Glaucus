package protocol

import "fmt"

type TLV struct {
	Tag    int
	Length int
	Value  string
}

func (tlv *TLV) String() string {
	return fmt.Sprintf("%04X%04X%s", tlv.Tag, tlv.Length, tlv.Value)
}

// TLV长度，包含头部
func (tlv *TLV) Size() int {
	return tlv.Length + 8
}