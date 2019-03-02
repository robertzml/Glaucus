package equipment

import (
	"fmt"
	"testing"
)

func TestReadWaterHeater(t *testing.T) {
	var w WaterHeater

	err := w.GetStatus("01100101801100e2")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", w)
}