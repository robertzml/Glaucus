package glog

import (
	"testing"
)

func TestWriteFile(t *testing.T) {
	_ := GlogPacket{Level: 1, PackageName: "test", Title: "testing", Message: "This is a sample."}

	//write(s)
}
