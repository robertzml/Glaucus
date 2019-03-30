package protocol

import (
	"fmt"
	"testing"
)

func TestTLV_String(t *testing.T) {
	tlv := TLV{10, 2, "abc"}

	// fmt.Println(tlv.String())

	if tlv.Tag != 10 {
		t.Error("error")
	}
}

func TestParseHead(t *testing.T) {
	msg := "Homeconsole02.00000000010003006500010006123456001B00080000003F0007000110004002458783926-533E-484B-9B79-FEE11E5A6832001A00013001C00010"

	seq, _, err := parseHead(msg)
	if err != nil {
		t.Error(err.Error())
	}

	if seq != "00000001" {
		t.Error("seq number is wrong")
	}
}

func TestParseMessage(t *testing.T) {
	msg := "Homeconsole02.00000000010003006500010006123456001B00080000003F0007000110004002458783926-533E-484B-9B79-FEE11E5A6832001A00013001C00010"

	_, payload, err := parseHead(msg)
	if err != nil {
		t.Error(err.Error())
	}

	tlv, err := parseTLV(payload, 0)
	if err != nil {
		t.Error(err.Error())
	}

	if tlv.Tag != 0x03 {
		t.Error("tag incorrect")
	}

	fmt.Printf("Cell length: %d\n", tlv.Length)

	index := 0
	for index < tlv.Length {
		item, err := parseTLV(tlv.Value, index);
		if err != nil {
			t.Error(err.Error())
			return
		}

		fmt.Printf("Tag: %d, Lenght: %d, Value: %s\n", item.Tag, item.Length, item.Value)
		index += item.Length + 8
	}
}
