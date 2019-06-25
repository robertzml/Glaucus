package tlv

import (
	"fmt"
	"testing"
)

func TestTimestamp(t *testing.T) {
	s := "0001010000"

	ts, err := ParseDateToTimestamp(s)

	if err == nil {
		fmt.Printf("ts: %d", ts)
	} else {
		fmt.Println(err)
	}
}