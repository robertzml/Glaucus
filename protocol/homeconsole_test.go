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
	msg := "Homeconsole05.00000000010003006500010006123456001B00080000003F0007000110004002458783926-533E-484B-9B79-FEE11E5A6832001A00013001C00010"

	seq, _, err := parseHead(msg)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if seq != "00000001" {
		t.Error("seq number is wrong")
	}
}

func TestParseMessage(t *testing.T) {
	msg := "Homeconsole05.000000028a0014007e0127001001100101801100e3012b000e1000000000006301250017hair salon water heater012a000bRH117.11.23012800160008000405000005000240"

	_, payload, err := parseHead(msg)
	if err != nil {
		t.Error(err.Error())
	}

	tlv, err := parseTLV(payload, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if tlv.Tag != 0x14 {
		t.Error("tag incorrect")
		return
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


//func TestReceive(t *testing.T) {
//	productType := 1
//	topic := "equipment/1/1/+/status_info"
//	payload := "Homeconsole05.000000028a0014007e0127001001100101801100e3012b000e1000000000006301250017hair salon water heater012a000bRH117.11.23012800160008000405000005000240"
//
//
//	base.InitConfig()
//	redis.InitPool(base.DefaultConfig.RedisDatabase)
//
//
//	Receive(productType, topic, []byte(payload), 2)
//}