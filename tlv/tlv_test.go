package tlv

import (
	"fmt"
	"testing"
)

func TestTimestamp(t *testing.T) {
	s := "0001010000"

	ts, err := ParseDateToTimestamp(s)

	if err == nil {
		fmt.Printf("ts: %d\n", ts)
	} else {
		fmt.Println(err)
	}
}

func TestTimestamp2(t *testing.T) {
	s := "0b1c1101e0"

	ts, err := ParseDateToTimestamp(s)

	if err == nil {
		fmt.Printf("ts: %d\n", ts)
	} else {
		fmt.Println(err)
	}
}

func TestTimestamp3(t *testing.T) {
	s := int64(946656000000)

	ts := GetCurDateTimeByTimestamp(s)

	fmt.Println(ts)
}
