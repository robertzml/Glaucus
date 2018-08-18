package protocol

import (
	"testing"
	"fmt"
)

func TestTLV_String(t *testing.T) {
	tlv:= TLV {10, 2, "abc"}

	fmt.Println(tlv.String())

	if tlv.Tag != 10 {
		t.Error("error")
	}
}

func TestParseHead(t *testing.T) {
	msg := "Homeconsole04.00000000010003006500010006123456001B00080000003F0007000110004002458783926-533E-484B-9B79-FEE11E5A6832001A00013001C00010"

	err := ParseHead(msg)
	if err != nil {
		t.Error(err.Error())
	}
}