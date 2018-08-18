package protocol

import (
	"fmt"
	"errors"
)

var version string = "Homeconsole04.00"

type TLV struct {
	Tag    int
	Length int
	Value  string
}

func (tlv *TLV) String() string {
	return fmt.Sprintf("%04X%04X%s", tlv.Tag, tlv.Length, tlv.Value)
}


type Message interface {
	Parse(input string, pos int) TLV
}

func ParseHead(message string) error{
	v := message[0: len(version)]
	if version != v {
		return errors.New("version not match")
	}

	// this.sequence = Convert.ToInt32(message.Substring(this.version.Length, 8), 16);
	seq := message[len(version): 16]
	fmt.Println(seq)

	return nil
}
