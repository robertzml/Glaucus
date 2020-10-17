package glog

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSerialize(t *testing.T) {
	pak := packet{Level: 1, System: "glaucus", Module: "glog", Action: "Test", Message: "test serialize"}
	jsonData, _ := json.Marshal(pak)
	fmt.Println(string(jsonData))
}